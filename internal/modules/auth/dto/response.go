package dto

import "github.com/open-suite/authorization/internal/entities"

type GoogleCallbackResponse struct {
	User    *entities.User `json:"user"`
	Google  GoogleUserInfo `json:"google"`
	Created bool           `json:"created"`
}

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture,omitempty"`
}

type SessionResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token,omitempty"`
	SessionState     string `json:"session_state,omitempty"`
	Scope            string `json:"scope,omitempty"`
}

type AccessSummaryResponse struct {
	App         string           `json:"app"`
	Menus       []AccessibleMenu `json:"menus"`
	Permissions []string         `json:"permissions"`
}

type AccessibleMenu struct {
	Code               string  `json:"code"`
	Path               string  `json:"path"`
	RequiredPermission *string `json:"required_permission"`
}

type UserAppAccessResponse struct {
	Items []UserAppAccess `json:"items"`
}

type UserAppAccess struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type AccessCheckResponse struct {
	Allowed    bool   `json:"allowed"`
	App        string `json:"app"`
	Permission string `json:"permission"`
}

type AccessTokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
	TokenType string `json:"token_type"`
}
