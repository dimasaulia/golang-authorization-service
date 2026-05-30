package services

import (
	"context"
	"encoding/json"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/auth/dto"
)

type AuthService interface {
	JWKS(ctx context.Context) (json.RawMessage, error)
	Login(ctx context.Context, request dto.LoginRequest) (*dto.SessionResponse, error)
	Refresh(ctx context.Context, request dto.RefreshRequest) (*dto.SessionResponse, error)
	Logout(ctx context.Context, request dto.LogoutRequest) error
	KeycloakRedirectURL(ctx context.Context, callbackURL string, prompt string) (string, error)
	HandleKeycloakCallback(ctx context.Context, code string, state string) (*dto.SessionResponse, string, error)
	HandleKeycloakErrorCallback(ctx context.Context, state string, reason string) (string, error)
	ExchangeKeycloakCallbackCode(ctx context.Context, code string) (*dto.SessionResponse, error)
	GoogleRedirectURL(ctx context.Context, organizationID int64) (string, error)
	HandleGoogleCallback(ctx context.Context, code string, state string) (*dto.GoogleCallbackResponse, string, error)
	CurrentUserAccess(ctx context.Context, userID int64) (*dto.CurrentUserAccessResponse, error)
	UpdateUser(ctx context.Context, userID int64, request dto.UpdateUserRequest) (*entities.User, error)
	AccessSummary(ctx context.Context, userID int64, appCode string) (*dto.AccessSummaryResponse, error)
	Apps(ctx context.Context, userID int64) (*dto.UserAppAccessResponse, error)
	Menus(ctx context.Context, userID int64, appCode string) ([]dto.AccessibleMenu, error)
	Permissions(ctx context.Context, userID int64, appCode string) ([]string, error)
	Check(ctx context.Context, userID int64, appCode string, permission string) (*dto.AccessCheckResponse, error)
	CheckPermission(ctx context.Context, userID int64, appCode string, permission string) (bool, error)
	AccessToken(ctx context.Context, userID int64, appCode string) (*dto.AccessTokenResponse, error)
}
