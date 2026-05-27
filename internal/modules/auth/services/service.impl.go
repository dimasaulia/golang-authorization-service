package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/auth/dto"
	"github.com/open-suite/authorization/internal/modules/auth/repositories"
	userRepositories "github.com/open-suite/authorization/internal/modules/users/repositories"
	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/platform/freeipa"
	"github.com/open-suite/authorization/internal/platform/keycloak"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/platform/redis"
	goredis "github.com/redis/go-redis/v9"
)

const (
	googleAuthorizeURL = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL     = "https://oauth2.googleapis.com/token"
	googleUserInfoURL  = "https://openidconnect.googleapis.com/v1/userinfo"
	googleProvider     = "google"
	freeIPAProvider    = "freeipa"
)

var (
	ErrGoogleOAuthNotConfigured = errors.New("google oauth is not configured")
	ErrInvalidOAuthState        = errors.New("invalid oauth state")
	ErrInvalidGoogleCallback    = errors.New("invalid google callback")
)

type AuthServiceImpl struct {
	cfg              config.Config
	redis            *redis.Redis
	userRepository   userRepositories.UserRepository
	accessRepository repositories.AccessRepository
	keycloak         keycloak.Client
	freeipa          freeipa.Client
	log              *logger.LayerLogger
	httpClient       *http.Client
}

type googleState struct {
	OrganizationID int64 `json:"organization_id"`
}

type googleTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	IDToken     string `json:"id_token"`
}

func NewAuthService(cfg config.Config, redisClient *redis.Redis, userRepository userRepositories.UserRepository, accessRepository repositories.AccessRepository, keycloakClient keycloak.Client, freeIPAClient freeipa.Client, appLogger *logger.Logger) AuthService {
	return &AuthServiceImpl{
		cfg:              cfg,
		redis:            redisClient,
		userRepository:   userRepository,
		accessRepository: accessRepository,
		keycloak:         keycloakClient,
		freeipa:          freeIPAClient,
		log:              appLogger.Layer("service.auth"),
		httpClient:       &http.Client{Timeout: 15 * time.Second},
	}
}

func (s *AuthServiceImpl) JWKS(ctx context.Context) (json.RawMessage, error) {
	end := s.log.Start(ctx, "JWKS")
	jwks, err := s.keycloak.JWKS(ctx)
	end(err)
	return jwks, err
}

func (s *AuthServiceImpl) Login(ctx context.Context, request dto.LoginRequest) (*dto.SessionResponse, error) {
	end := s.log.Start(ctx, "Login")
	token, err := s.keycloak.Login(ctx, keycloak.LoginInput{
		Username: strings.TrimSpace(request.Username),
		Password: request.Password,
	})
	if err != nil {
		end(err)
		return nil, err
	}

	end(nil)
	return sessionResponseFromToken(token), nil
}

func (s *AuthServiceImpl) Refresh(ctx context.Context, request dto.RefreshRequest) (*dto.SessionResponse, error) {
	end := s.log.Start(ctx, "Refresh")
	token, err := s.keycloak.Refresh(ctx, keycloak.RefreshInput{
		RefreshToken: strings.TrimSpace(request.RefreshToken),
	})
	if err != nil {
		end(err)
		return nil, err
	}

	end(nil)
	return sessionResponseFromToken(token), nil
}

func (s *AuthServiceImpl) Logout(ctx context.Context, request dto.LogoutRequest) error {
	end := s.log.Start(ctx, "Logout")
	err := s.keycloak.Logout(ctx, keycloak.LogoutInput{
		RefreshToken: strings.TrimSpace(request.RefreshToken),
	})
	end(err)
	return err
}

func (s *AuthServiceImpl) GoogleRedirectURL(ctx context.Context, organizationID int64) (string, error) {
	end := s.log.Start(ctx, "GoogleRedirectURL")
	if err := s.validateGoogleConfig(); err != nil {
		end(err)
		return "", err
	}
	if organizationID == 0 {
		organizationID = s.cfg.OAuth.Google.DefaultOrganizationID
	}
	if organizationID == 0 {
		err := ErrInvalidGoogleCallback
		end(err)
		return "", err
	}

	state := newStateToken()
	payload, err := json.Marshal(googleState{OrganizationID: organizationID})
	if err != nil {
		end(err)
		return "", err
	}

	if err := s.redis.Client.Set(ctx, googleStateKey(state), payload, 10*time.Minute).Err(); err != nil {
		end(err)
		return "", err
	}

	values := url.Values{}
	values.Set("client_id", s.cfg.OAuth.Google.ClientID)
	values.Set("redirect_uri", s.cfg.OAuth.Google.RedirectURL)
	values.Set("response_type", "code")
	values.Set("scope", strings.Join(s.cfg.OAuth.Google.Scopes, " "))
	values.Set("state", state)
	values.Set("access_type", "offline")
	values.Set("prompt", "select_account")

	redirectURL := googleAuthorizeURL + "?" + values.Encode()
	end(nil)
	return redirectURL, nil
}

func (s *AuthServiceImpl) HandleGoogleCallback(ctx context.Context, code string, state string) (*dto.GoogleCallbackResponse, string, error) {
	end := s.log.Start(ctx, "HandleGoogleCallback")
	code = strings.TrimSpace(code)
	state = strings.TrimSpace(state)
	if code == "" || state == "" {
		end(ErrInvalidGoogleCallback)
		return nil, s.failureRedirectURL("invalid_callback"), ErrInvalidGoogleCallback
	}

	storedState, err := s.consumeGoogleState(ctx, state)
	if err != nil {
		end(err)
		return nil, s.failureRedirectURL("invalid_state"), err
	}

	token, err := s.exchangeCode(ctx, code)
	if err != nil {
		end(err)
		return nil, s.failureRedirectURL("token_exchange_failed"), err
	}

	googleUser, err := s.fetchGoogleUser(ctx, token.AccessToken)
	if err != nil {
		end(err)
		return nil, s.failureRedirectURL("userinfo_failed"), err
	}
	if googleUser.Sub == "" || googleUser.Email == "" {
		end(ErrInvalidGoogleCallback)
		return nil, s.failureRedirectURL("invalid_userinfo"), ErrInvalidGoogleCallback
	}

	user, created, err := s.findOrCreateGoogleUser(ctx, storedState.OrganizationID, googleUser)
	if err != nil {
		end(err)
		return nil, s.failureRedirectURL("user_sync_failed"), err
	}

	redirectURL := s.successRedirectURL(user.ID)
	end(nil, "user_id", user.ID, "created", created)
	return &dto.GoogleCallbackResponse{
		User:    user,
		Google:  googleUser,
		Created: created,
	}, redirectURL, nil
}

func (s *AuthServiceImpl) validateGoogleConfig() error {
	if s.cfg.OAuth.Google.ClientID == "" || s.cfg.OAuth.Google.ClientSecret == "" || s.cfg.OAuth.Google.RedirectURL == "" {
		return ErrGoogleOAuthNotConfigured
	}
	return nil
}

func (s *AuthServiceImpl) AccessSummary(ctx context.Context, userID int64, appCode string) (*dto.AccessSummaryResponse, error) {
	appCode = strings.TrimSpace(appCode)
	cacheKey := accessCacheKey(userID, appCode)
	var cached dto.AccessSummaryResponse
	if ok := s.getCache(ctx, cacheKey, &cached); ok {
		return &cached, nil
	}

	item, err := s.accessRepository.FindAccessSummary(ctx, userID, appCode)
	if err != nil {
		return nil, err
	}
	s.setCache(ctx, cacheKey, item)
	return item, nil
}

func (s *AuthServiceImpl) Apps(ctx context.Context, userID int64) (*dto.UserAppAccessResponse, error) {
	cacheKey := appsCacheKey(userID)
	var cached dto.UserAppAccessResponse
	if ok := s.getCache(ctx, cacheKey, &cached); ok {
		return &cached, nil
	}

	items, err := s.accessRepository.FindApps(ctx, userID)
	if err != nil {
		return nil, err
	}
	response := &dto.UserAppAccessResponse{Items: items}
	s.setCache(ctx, cacheKey, response)
	return response, nil
}

func (s *AuthServiceImpl) Menus(ctx context.Context, userID int64, appCode string) ([]dto.AccessibleMenu, error) {
	summary, err := s.AccessSummary(ctx, userID, appCode)
	if err != nil {
		return nil, err
	}
	return summary.Menus, nil
}

func (s *AuthServiceImpl) Permissions(ctx context.Context, userID int64, appCode string) ([]string, error) {
	summary, err := s.AccessSummary(ctx, userID, appCode)
	if err != nil {
		return nil, err
	}
	return summary.Permissions, nil
}

func (s *AuthServiceImpl) Check(ctx context.Context, userID int64, appCode string, permission string) (*dto.AccessCheckResponse, error) {
	permission = strings.TrimSpace(permission)
	summary, err := s.AccessSummary(ctx, userID, appCode)
	if err != nil {
		return nil, err
	}
	allowed := false
	for _, item := range summary.Permissions {
		if strings.EqualFold(item, permission) {
			allowed = true
			break
		}
	}
	return &dto.AccessCheckResponse{
		Allowed:    allowed,
		App:        summary.App,
		Permission: permission,
	}, nil
}

func (s *AuthServiceImpl) CheckPermission(ctx context.Context, userID int64, appCode string, permission string) (bool, error) {
	result, err := s.Check(ctx, userID, appCode, permission)
	if err != nil {
		return false, err
	}
	return result.Allowed, nil
}

func (s *AuthServiceImpl) AccessToken(ctx context.Context, userID int64, appCode string) (*dto.AccessTokenResponse, error) {
	summary, err := s.AccessSummary(ctx, userID, appCode)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(s.cfg.Authz.TokenTTL)
	token, err := s.signAccessToken(userID, summary, expiresAt)
	if err != nil {
		return nil, err
	}
	return &dto.AccessTokenResponse{
		Token:     token,
		ExpiresIn: int64(time.Until(expiresAt).Seconds()),
		TokenType: "Bearer",
	}, nil
}

func (s *AuthServiceImpl) consumeGoogleState(ctx context.Context, state string) (googleState, error) {
	key := googleStateKey(state)
	raw, err := s.redis.Client.Get(ctx, key).Bytes()
	if errors.Is(err, goredis.Nil) {
		return googleState{}, ErrInvalidOAuthState
	}
	if err != nil {
		return googleState{}, err
	}
	if err := s.redis.Client.Del(ctx, key).Err(); err != nil {
		return googleState{}, err
	}

	var payload googleState
	if err := json.Unmarshal(raw, &payload); err != nil {
		return googleState{}, err
	}
	if payload.OrganizationID == 0 {
		return googleState{}, ErrInvalidOAuthState
	}
	return payload, nil
}

func (s *AuthServiceImpl) exchangeCode(ctx context.Context, code string) (googleTokenResponse, error) {
	values := url.Values{}
	values.Set("code", code)
	values.Set("client_id", s.cfg.OAuth.Google.ClientID)
	values.Set("client_secret", s.cfg.OAuth.Google.ClientSecret)
	values.Set("redirect_uri", s.cfg.OAuth.Google.RedirectURL)
	values.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, googleTokenURL, bytes.NewBufferString(values.Encode()))
	if err != nil {
		return googleTokenResponse{}, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var token googleTokenResponse
	if err := s.doJSON(req, &token); err != nil {
		return googleTokenResponse{}, err
	}
	if token.AccessToken == "" {
		return googleTokenResponse{}, ErrInvalidGoogleCallback
	}
	return token, nil
}

func (s *AuthServiceImpl) fetchGoogleUser(ctx context.Context, accessToken string) (dto.GoogleUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, googleUserInfoURL, nil)
	if err != nil {
		return dto.GoogleUserInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	var user dto.GoogleUserInfo
	if err := s.doJSON(req, &user); err != nil {
		return dto.GoogleUserInfo{}, err
	}
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Name = strings.TrimSpace(user.Name)
	return user, nil
}

func (s *AuthServiceImpl) doJSON(req *http.Request, target any) error {
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("google oauth request failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	return json.Unmarshal(body, target)
}

func (s *AuthServiceImpl) findOrCreateGoogleUser(ctx context.Context, organizationID int64, googleUser dto.GoogleUserInfo) (*entities.User, bool, error) {
	user, err := s.userRepository.FindByIdentity(ctx, googleProvider, googleUser.Sub)
	if err == nil {
		if googleUser.EmailVerified && user.EmailVerifiedAt == nil {
			_ = s.userRepository.MarkEmailVerified(ctx, user.ID)
			user, _ = s.userRepository.FindByID(ctx, user.ID)
		}
		if err := s.ensureGoogleUserInFreeIPA(ctx, user, false); err != nil {
			return nil, false, err
		}
		return user, false, nil
	}
	if !errors.Is(err, userRepositories.ErrNotFound) {
		return nil, false, err
	}

	username := usernameFromGoogleUser(googleUser)
	email := googleUser.Email
	existingUser, err := s.userRepository.FindByEmail(ctx, email)
	if err == nil {
		if err := s.userRepository.LinkIdentity(ctx, entities.UserIdentity{
			UserId:         existingUser.ID,
			Provider:       googleProvider,
			ProviderUserId: googleUser.Sub,
			Username:       &username,
			Email:          &email,
			IsPrimary:      false,
		}); err != nil {
			return nil, false, err
		}
		if googleUser.EmailVerified && existingUser.EmailVerifiedAt == nil {
			_ = s.userRepository.MarkEmailVerified(ctx, existingUser.ID)
			existingUser, _ = s.userRepository.FindByID(ctx, existingUser.ID)
		}
		if err := s.ensureGoogleUserInFreeIPA(ctx, existingUser, false); err != nil {
			return nil, false, err
		}
		return existingUser, false, nil
	}
	if !errors.Is(err, userRepositories.ErrNotFound) {
		return nil, false, err
	}

	displayName := googleUser.Name
	if displayName == "" {
		displayName = email
	}
	user, err = s.userRepository.Create(ctx, userRepositories.CreateUserInput{
		User: entities.User{
			OrganizationId: organizationID,
			Username:       username,
			Email:          email,
			DisplayName:    displayName,
			Type:           "external",
			Status:         "active",
		},
		Identity: &entities.UserIdentity{
			Provider:       googleProvider,
			ProviderUserId: googleUser.Sub,
			Username:       &username,
			Email:          &email,
			IsPrimary:      true,
		},
	})
	if err != nil {
		return nil, false, err
	}
	if googleUser.EmailVerified {
		_ = s.userRepository.MarkEmailVerified(ctx, user.ID)
		user, _ = s.userRepository.FindByID(ctx, user.ID)
	}
	if err := s.ensureGoogleUserInFreeIPA(ctx, user, true); err != nil {
		_ = s.userRepository.Delete(ctx, user.ID)
		return nil, false, err
	}
	return user, true, nil
}

func (s *AuthServiceImpl) ensureGoogleUserInFreeIPA(ctx context.Context, user *entities.User, rollbackFreeIPAOnLinkFailure bool) error {
	uid, err := s.freeipa.CreateUser(ctx, freeipa.CreateUserInput{
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
	})
	if err != nil {
		return err
	}

	if existingUser, err := s.userRepository.FindByIdentity(ctx, freeIPAProvider, uid); err == nil && existingUser.ID == user.ID {
		return nil
	}

	username := user.Username
	email := user.Email
	if err := s.userRepository.LinkIdentity(ctx, entities.UserIdentity{
		UserId:         user.ID,
		Provider:       freeIPAProvider,
		ProviderUserId: uid,
		Username:       &username,
		Email:          &email,
		IsPrimary:      false,
	}); err != nil {
		if rollbackFreeIPAOnLinkFailure {
			_ = s.freeipa.DeleteUser(ctx, uid)
		}
		return err
	}

	return nil
}

func usernameFromGoogleUser(user dto.GoogleUserInfo) string {
	local := strings.Split(user.Email, "@")[0]
	local = strings.TrimSpace(local)
	if local == "" {
		local = "google_user"
	}
	sub := strings.TrimSpace(user.Sub)
	if len(sub) > 8 {
		sub = sub[:8]
	}
	if sub == "" {
		return local
	}
	return local + "_" + sub
}

func (s *AuthServiceImpl) successRedirectURL(userID int64) string {
	if s.cfg.OAuth.Google.SuccessRedirectURL == "" {
		return ""
	}
	values := url.Values{}
	values.Set("provider", googleProvider)
	values.Set("user_id", strconv.FormatInt(userID, 10))
	return s.cfg.OAuth.Google.SuccessRedirectURL + "?" + values.Encode()
}

func (s *AuthServiceImpl) failureRedirectURL(reason string) string {
	if s.cfg.OAuth.Google.FailureRedirectURL == "" {
		return ""
	}
	values := url.Values{}
	values.Set("provider", googleProvider)
	values.Set("reason", reason)
	return s.cfg.OAuth.Google.FailureRedirectURL + "?" + values.Encode()
}

func googleStateKey(state string) string {
	return "oauth:google:state:" + state
}

func newStateToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

func (s *AuthServiceImpl) getCache(ctx context.Context, key string, target any) bool {
	raw, err := s.redis.Client.Get(ctx, key).Bytes()
	if err != nil {
		return false
	}
	return json.Unmarshal(raw, target) == nil
}

func (s *AuthServiceImpl) setCache(ctx context.Context, key string, value any) {
	raw, err := json.Marshal(value)
	if err != nil {
		return
	}
	if err := s.redis.Client.Set(ctx, key, raw, s.cfg.Authz.CacheTTL).Err(); err != nil {
		s.log.Warn(ctx, "cache.set_failed", "key", key, "error", err.Error())
	}
}

func (s *AuthServiceImpl) signAccessToken(userID int64, summary *dto.AccessSummaryResponse, expiresAt time.Time) (string, error) {
	header := map[string]any{
		"alg": "HS256",
		"typ": "JWT",
	}
	payload := map[string]any{
		"sub":         fmt.Sprintf("%d", userID),
		"app":         summary.App,
		"menus":       summary.Menus,
		"permissions": summary.Permissions,
		"exp":         expiresAt.Unix(),
		"iat":         time.Now().Unix(),
	}

	headerPart, err := encodeJWTPart(header)
	if err != nil {
		return "", err
	}
	payloadPart, err := encodeJWTPart(payload)
	if err != nil {
		return "", err
	}
	unsigned := headerPart + "." + payloadPart
	signature := hmac.New(sha256.New, []byte(s.cfg.Authz.TokenSecret))
	signature.Write([]byte(unsigned))
	return unsigned + "." + base64.RawURLEncoding.EncodeToString(signature.Sum(nil)), nil
}

func encodeJWTPart(value any) (string, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func accessCacheKey(userID int64, appCode string) string {
	return fmt.Sprintf("authz:access:user:%d:app:%s", userID, strings.ToLower(strings.TrimSpace(appCode)))
}

func appsCacheKey(userID int64) string {
	return fmt.Sprintf("authz:apps:user:%d", userID)
}

func sessionResponseFromToken(token *keycloak.TokenSet) *dto.SessionResponse {
	return &dto.SessionResponse{
		AccessToken:      token.AccessToken,
		ExpiresIn:        token.ExpiresIn,
		RefreshExpiresIn: token.RefreshExpiresIn,
		RefreshToken:     token.RefreshToken,
		TokenType:        token.TokenType,
		IDToken:          token.IDToken,
		SessionState:     token.SessionState,
		Scope:            token.Scope,
	}
}
