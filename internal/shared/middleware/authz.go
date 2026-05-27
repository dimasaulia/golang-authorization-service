package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/platform/database"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/platform/redis"
	"github.com/open-suite/authorization/internal/shared/requestctx"
	"github.com/open-suite/authorization/internal/shared/response"
)

const (
	keycloakProvider = "keycloak"
	userCacheTTL     = 15 * time.Minute
	jwksCacheTTL     = 15 * time.Minute
)

var (
	errInvalidToken = errors.New("invalid bearer token")
)

type Authenticator struct {
	cfg        config.Config
	db         *database.Database
	redis      *redis.Redis
	checker    PermissionChecker
	response   *response.Sender
	log        *logger.LayerLogger
	httpClient *http.Client

	jwksMu        sync.RWMutex
	jwksExpiresAt time.Time
	publicKeys    map[string]*rsa.PublicKey
}

type PermissionChecker interface {
	CheckPermission(ctx context.Context, userID int64, appCode string, permission string) (bool, error)
}

type jwksResponse struct {
	Keys []jwkKey `json:"keys"`
}

type jwkKey struct {
	KID string `json:"kid"`
	KTY string `json:"kty"`
	ALG string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type authClaims struct {
	Subject           string `json:"sub"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthenticator(cfg config.Config, db *database.Database, redisClient *redis.Redis, checker PermissionChecker, sender *response.Sender, appLogger *logger.Logger) *Authenticator {
	return &Authenticator{
		cfg:        cfg,
		db:         db,
		redis:      redisClient,
		checker:    checker,
		response:   sender,
		log:        appLogger.Layer("middleware.authz"),
		httpClient: &http.Client{Timeout: 10 * time.Second},
		publicKeys: map[string]*rsa.PublicKey{},
	}
}

func (a *Authenticator) AuthenticateOptional(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r)
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		claims, userID, err := a.authenticate(r.Context(), token)
		if err != nil {
			a.log.Warn(r.Context(), "authenticate.failed", "error", err.Error())
			a.response.Error(w, r, http.StatusUnauthorized, "auth.unauthorized", nil)
			return
		}

		ctx := requestctx.WithSubject(r.Context(), claims.Subject)
		ctx = requestctx.WithUserID(ctx, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Authenticator) RequireAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestctx.UserID(r.Context()) > 0 {
			next.ServeHTTP(w, r)
			return
		}

		token := bearerToken(r)
		if token == "" {
			a.response.Error(w, r, http.StatusUnauthorized, "auth.unauthorized", nil)
			return
		}

		claims, userID, err := a.authenticate(r.Context(), token)
		if err != nil {
			a.log.Warn(r.Context(), "authenticate.failed", "error", err.Error())
			a.response.Error(w, r, http.StatusUnauthorized, "auth.unauthorized", nil)
			return
		}

		ctx := requestctx.WithSubject(r.Context(), claims.Subject)
		ctx = requestctx.WithUserID(ctx, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Authenticator) RequirePermission(appCode string, permission string) Middleware {
	return func(next http.Handler) http.Handler {
		return a.RequireAuthenticated(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := requestctx.UserID(r.Context())
			allowed, err := a.checker.CheckPermission(r.Context(), userID, appCode, permission)
			if err != nil {
				a.log.Warn(r.Context(), "permission.check_failed", "error", err.Error(), "user_id", userID, "app", appCode, "permission", permission)
				a.response.Error(w, r, http.StatusForbidden, "auth.forbidden", nil)
				return
			}
			if !allowed {
				a.response.Error(w, r, http.StatusForbidden, "auth.forbidden", nil)
				return
			}

			next.ServeHTTP(w, r)
		}))
	}
}

func (a *Authenticator) authenticate(ctx context.Context, rawToken string) (*authClaims, int64, error) {
	claims, err := a.verifyToken(ctx, rawToken)
	if err != nil {
		return nil, 0, err
	}
	if claims.Subject == "" {
		return nil, 0, errInvalidToken
	}

	userID, err := a.resolveUserID(ctx, claims.Subject, claims.PreferredUsername, claims.Email)
	if err != nil {
		return nil, 0, err
	}
	return claims, userID, nil
}

func (a *Authenticator) verifyToken(ctx context.Context, rawToken string) (*authClaims, error) {
	claims := &authClaims{}
	token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errInvalidToken
		}
		kid, _ := token.Header["kid"].(string)
		if kid == "" {
			return nil, errInvalidToken
		}
		return a.publicKey(ctx, kid)
	})
	if err != nil || token == nil || !token.Valid {
		return nil, errInvalidToken
	}

	if claims.Issuer != a.keycloakIssuer() {
		return nil, errInvalidToken
	}
	return claims, nil
}

func (a *Authenticator) publicKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	if key := a.cachedPublicKey(kid); key != nil {
		return key, nil
	}
	if err := a.refreshJWKS(ctx); err != nil {
		return nil, err
	}
	if key := a.cachedPublicKey(kid); key != nil {
		return key, nil
	}
	return nil, errInvalidToken
}

func (a *Authenticator) cachedPublicKey(kid string) *rsa.PublicKey {
	a.jwksMu.RLock()
	defer a.jwksMu.RUnlock()
	if time.Now().After(a.jwksExpiresAt) {
		return nil
	}
	return a.publicKeys[kid]
}

func (a *Authenticator) refreshJWKS(ctx context.Context) error {
	if a.cfg.Keycloak.BaseURL == "" || a.cfg.Keycloak.Realm == "" {
		return errInvalidToken
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.keycloakIssuer()+"/protocol/openid-connect/certs", nil)
	if err != nil {
		return err
	}
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("jwks request failed: status=%d", resp.StatusCode)
	}

	var jwks jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}
	keys := make(map[string]*rsa.PublicKey, len(jwks.Keys))
	for _, item := range jwks.Keys {
		key, err := rsaPublicKey(item)
		if err != nil {
			continue
		}
		keys[item.KID] = key
	}
	if len(keys) == 0 {
		return errInvalidToken
	}

	a.jwksMu.Lock()
	a.publicKeys = keys
	a.jwksExpiresAt = time.Now().Add(jwksCacheTTL)
	a.jwksMu.Unlock()
	return nil
}

func (a *Authenticator) resolveUserID(ctx context.Context, subject string, username string, email string) (int64, error) {
	key := userIdentityCacheKey(subject)
	if cached, err := a.redis.Client.Get(ctx, key).Int64(); err == nil && cached > 0 {
		return cached, nil
	}

	query := `
		SELECT id
		FROM (
			SELECT u.id, 1 AS priority
			FROM user_identities ui
			JOIN users u ON u.id = ui.user_id
			WHERE ui.provider = $1
			AND ui.provider_user_id = $2
			AND u.status = 'active'
			UNION ALL
			SELECT u.id, 2 AS priority
			FROM users u
			WHERE u.status = 'active'
			AND (
				($3 <> '' AND LOWER(u.username) = LOWER($3))
				OR ($4 <> '' AND LOWER(u.email) = LOWER($4))
			)
		) found
		ORDER BY priority ASC
		LIMIT 1`

	var userID int64
	if err := a.db.Pool.QueryRow(ctx, query, keycloakProvider, subject, strings.TrimSpace(username), strings.TrimSpace(email)).Scan(&userID); err != nil {
		return 0, err
	}

	_ = a.redis.Client.Set(ctx, key, userID, userCacheTTL).Err()
	return userID, nil
}

func (a *Authenticator) keycloakIssuer() string {
	return strings.TrimRight(a.cfg.Keycloak.BaseURL, "/") + "/realms/" + a.cfg.Keycloak.Realm
}

func bearerToken(r *http.Request) string {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func rsaPublicKey(key jwkKey) (*rsa.PublicKey, error) {
	if key.KTY != "RSA" || key.N == "" || key.E == "" || key.KID == "" {
		return nil, errInvalidToken
	}
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, err
	}
	exponent := 0
	for _, b := range eBytes {
		exponent = exponent<<8 + int(b)
	}
	if exponent == 0 {
		return nil, errInvalidToken
	}
	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: exponent,
	}, nil
}

func userIdentityCacheKey(subject string) string {
	return "sdk:authn:keycloak:sub:" + subject + ":user_id"
}
