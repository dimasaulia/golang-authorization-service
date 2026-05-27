package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/platform/logger"
)

const passwordCredentialType = "password"

var (
	ErrDisabled      = errors.New("keycloak integration is disabled")
	ErrNotConfigured = errors.New("keycloak integration is not configured")
)

type Client interface {
	Enabled() bool
	JWKS(ctx context.Context) (json.RawMessage, error)
	Login(ctx context.Context, input LoginInput) (*TokenSet, error)
	Refresh(ctx context.Context, input RefreshInput) (*TokenSet, error)
	Logout(ctx context.Context, input LogoutInput) error
	CreateUser(ctx context.Context, input CreateUserInput) (string, error)
	DeleteUser(ctx context.Context, userID string) error
}

type Keycloak struct {
	cfg        config.KeycloakConfig
	log        *logger.LayerLogger
	httpClient *http.Client
}

type CreateUserInput struct {
	Username            string
	Email               string
	DisplayName         string
	Enabled             bool
	EmailVerified       bool
	Password            string
	TemporaryPassword   bool
	FederatedIdentities []FederatedIdentity
}

type FederatedIdentity struct {
	IdentityProvider string
	UserID           string
	UserName         string
}

type LoginInput struct {
	Username string
	Password string
	Scope    string
}

type RefreshInput struct {
	RefreshToken string
}

type LogoutInput struct {
	RefreshToken string
}

type TokenSet struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token,omitempty"`
	NotBeforePolicy  int64  `json:"not-before-policy,omitempty"`
	SessionState     string `json:"session_state,omitempty"`
	Scope            string `json:"scope,omitempty"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

type userRepresentation struct {
	ID                  string                            `json:"id,omitempty"`
	Username            string                            `json:"username"`
	Email               string                            `json:"email"`
	FirstName           string                            `json:"firstName,omitempty"`
	LastName            string                            `json:"lastName,omitempty"`
	Enabled             bool                              `json:"enabled"`
	EmailVerified       bool                              `json:"emailVerified"`
	Credentials         []credentialRepresentation        `json:"credentials,omitempty"`
	FederatedIdentities []federatedIdentityRepresentation `json:"federatedIdentities,omitempty"`
}

type credentialRepresentation struct {
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

type federatedIdentityRepresentation struct {
	IdentityProvider string `json:"identityProvider"`
	UserID           string `json:"userId"`
	UserName         string `json:"userName"`
}

func New(cfg config.Config, appLogger *logger.Logger) Client {
	return &Keycloak{
		cfg:        cfg.Keycloak,
		log:        appLogger.Layer("platform.keycloak"),
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (k *Keycloak) Enabled() bool {
	return k.cfg.Enabled
}

func (k *Keycloak) JWKS(ctx context.Context) (json.RawMessage, error) {
	end := k.log.Start(ctx, "JWKS")
	if !k.cfg.Enabled {
		end(ErrDisabled)
		return nil, ErrDisabled
	}
	if k.cfg.BaseURL == "" || k.cfg.Realm == "" {
		end(ErrNotConfigured)
		return nil, ErrNotConfigured
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, k.realmURL("/protocol/openid-connect/certs"), nil)
	if err != nil {
		end(err)
		return nil, err
	}
	resp, err := k.httpClient.Do(req)
	if err != nil {
		end(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		end(err)
		return nil, err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		err := fmt.Errorf("keycloak jwks failed: status=%d body=%s", resp.StatusCode, string(body))
		end(err)
		return nil, err
	}

	end(nil)
	return json.RawMessage(body), nil
}

func (k *Keycloak) Login(ctx context.Context, input LoginInput) (*TokenSet, error) {
	end := k.log.Start(ctx, "Login")
	if !k.cfg.Enabled {
		end(ErrDisabled)
		return nil, ErrDisabled
	}
	if err := k.validateLoginConfig(); err != nil {
		end(err)
		return nil, err
	}

	values := k.loginClientValues()
	values.Set("grant_type", "password")
	values.Set("username", input.Username)
	values.Set("password", input.Password)
	values.Set("scope", valueOrDefault(input.Scope, "openid profile email"))

	token, err := k.tokenRequest(ctx, values)
	end(err)
	return token, err
}

func (k *Keycloak) Refresh(ctx context.Context, input RefreshInput) (*TokenSet, error) {
	end := k.log.Start(ctx, "Refresh")
	if !k.cfg.Enabled {
		end(ErrDisabled)
		return nil, ErrDisabled
	}
	if err := k.validateLoginConfig(); err != nil {
		end(err)
		return nil, err
	}

	values := k.loginClientValues()
	values.Set("grant_type", "refresh_token")
	values.Set("refresh_token", input.RefreshToken)

	token, err := k.tokenRequest(ctx, values)
	end(err)
	return token, err
}

func (k *Keycloak) Logout(ctx context.Context, input LogoutInput) error {
	end := k.log.Start(ctx, "Logout")
	if !k.cfg.Enabled {
		end(ErrDisabled)
		return ErrDisabled
	}
	if err := k.validateLoginConfig(); err != nil {
		end(err)
		return err
	}

	values := k.loginClientValues()
	values.Set("refresh_token", input.RefreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, k.realmURL("/protocol/openid-connect/logout"), strings.NewReader(values.Encode()))
	if err != nil {
		end(err)
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := k.httpClient.Do(req)
	if err != nil {
		end(err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		err := responseError("keycloak logout failed", resp)
		end(err)
		return err
	}

	end(nil)
	return nil
}

func (k *Keycloak) CreateUser(ctx context.Context, input CreateUserInput) (string, error) {
	end := k.log.Start(ctx, "CreateUser")
	if !k.cfg.Enabled {
		end(ErrDisabled)
		return "", ErrDisabled
	}
	if err := k.validateConfig(); err != nil {
		end(err)
		return "", err
	}

	token, err := k.adminToken(ctx)
	if err != nil {
		end(err)
		return "", err
	}

	body := userRepresentation{
		Username:      strings.TrimSpace(input.Username),
		Email:         strings.ToLower(strings.TrimSpace(input.Email)),
		Enabled:       input.Enabled,
		EmailVerified: input.EmailVerified,
	}
	body.FirstName, body.LastName = splitDisplayName(input.DisplayName)
	if input.Password != "" {
		body.Credentials = []credentialRepresentation{{
			Type:      passwordCredentialType,
			Value:     input.Password,
			Temporary: input.TemporaryPassword,
		}}
	}
	for _, identity := range input.FederatedIdentities {
		body.FederatedIdentities = append(body.FederatedIdentities, federatedIdentityRepresentation{
			IdentityProvider: identity.IdentityProvider,
			UserID:           identity.UserID,
			UserName:         identity.UserName,
		})
	}

	payload, err := json.Marshal(body)
	if err != nil {
		end(err)
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, k.adminURL("/users"), bytes.NewReader(payload))
	if err != nil {
		end(err)
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := k.httpClient.Do(req)
	if err != nil {
		end(err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		userID, findErr := k.findUserID(ctx, token, body.Username, body.Email)
		if findErr != nil {
			end(findErr)
			return "", findErr
		}
		end(nil, "user_id", userID, "conflict", true)
		return userID, nil
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		err := responseError("create keycloak user failed", resp)
		end(err)
		return "", err
	}

	userID := idFromLocation(resp.Header.Get("Location"))
	if userID == "" {
		userID, err = k.findUserID(ctx, token, body.Username, body.Email)
		if err != nil {
			end(err)
			return "", err
		}
	}

	end(nil, "user_id", userID)
	return userID, nil
}

func (k *Keycloak) DeleteUser(ctx context.Context, userID string) error {
	end := k.log.Start(ctx, "DeleteUser", "user_id", userID)
	if !k.cfg.Enabled {
		end(ErrDisabled)
		return ErrDisabled
	}
	if userID == "" {
		end(nil)
		return nil
	}
	token, err := k.adminToken(ctx)
	if err != nil {
		end(err)
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, k.adminURL("/users/"+url.PathEscape(userID)), nil)
	if err != nil {
		end(err)
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := k.httpClient.Do(req)
	if err != nil {
		end(err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		err := responseError("delete keycloak user failed", resp)
		end(err)
		return err
	}

	end(nil)
	return nil
}

func (k *Keycloak) validateConfig() error {
	if k.cfg.BaseURL == "" || k.cfg.Realm == "" || k.cfg.ClientID == "" || k.cfg.ClientSecret == "" {
		return ErrNotConfigured
	}
	return nil
}

func (k *Keycloak) validateLoginConfig() error {
	if k.cfg.BaseURL == "" || k.cfg.Realm == "" || k.cfg.LoginClientID == "" {
		return ErrNotConfigured
	}
	return nil
}

func (k *Keycloak) loginClientValues() url.Values {
	values := url.Values{}
	values.Set("client_id", k.cfg.LoginClientID)
	if k.cfg.LoginClientSecret != "" {
		values.Set("client_secret", k.cfg.LoginClientSecret)
	}
	return values
}

func (k *Keycloak) tokenRequest(ctx context.Context, values url.Values) (*TokenSet, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, k.realmURL("/protocol/openid-connect/token"), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var token TokenSet
	if err := k.doJSON(req, &token); err != nil {
		return nil, err
	}
	if token.AccessToken == "" {
		return nil, errors.New("keycloak access token is empty")
	}
	return &token, nil
}

func (k *Keycloak) adminToken(ctx context.Context) (string, error) {
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", k.cfg.ClientID)
	values.Set("client_secret", k.cfg.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, k.realmURL("/protocol/openid-connect/token"), strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var token tokenResponse
	if err := k.doJSON(req, &token); err != nil {
		return "", err
	}
	if token.AccessToken == "" {
		return "", errors.New("keycloak admin token is empty")
	}
	return token.AccessToken, nil
}

func (k *Keycloak) findUserID(ctx context.Context, token string, username string, email string) (string, error) {
	values := url.Values{}
	values.Set("username", username)
	values.Set("exact", "true")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, k.adminURL("/users")+"?"+values.Encode(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	var users []userRepresentation
	if err := k.doJSON(req, &users); err != nil {
		return "", err
	}
	for _, user := range users {
		if strings.EqualFold(user.Username, username) || strings.EqualFold(user.Email, email) {
			return user.ID, nil
		}
	}
	return "", errors.New("keycloak user id not found")
}

func (k *Keycloak) doJSON(req *http.Request, target any) error {
	resp, err := k.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("keycloak request failed: status=%d body=%s", resp.StatusCode, string(body))
	}

	return json.Unmarshal(body, target)
}

func (k *Keycloak) realmURL(path string) string {
	return k.cfg.BaseURL + "/realms/" + url.PathEscape(k.cfg.Realm) + path
}

func (k *Keycloak) adminURL(path string) string {
	return k.cfg.BaseURL + "/admin/realms/" + url.PathEscape(k.cfg.Realm) + path
}

func idFromLocation(location string) string {
	if location == "" {
		return ""
	}
	parts := strings.Split(strings.TrimRight(location, "/"), "/")
	return parts[len(parts)-1]
}

func splitDisplayName(displayName string) (string, string) {
	parts := strings.Fields(displayName)
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

func responseError(message string, resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	return fmt.Errorf("%s: status=%d body=%s", message, resp.StatusCode, string(body))
}

func valueOrDefault(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
