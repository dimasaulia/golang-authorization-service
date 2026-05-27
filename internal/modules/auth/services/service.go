package services

import (
	"context"

	"github.com/open-suite/authorization/internal/modules/auth/dto"
)

type AuthService interface {
	GoogleRedirectURL(ctx context.Context, organizationID int64) (string, error)
	HandleGoogleCallback(ctx context.Context, code string, state string) (*dto.GoogleCallbackResponse, string, error)
	AccessSummary(ctx context.Context, userID int64, appCode string) (*dto.AccessSummaryResponse, error)
	Apps(ctx context.Context, userID int64) (*dto.UserAppAccessResponse, error)
	Menus(ctx context.Context, userID int64, appCode string) ([]dto.AccessibleMenu, error)
	Permissions(ctx context.Context, userID int64, appCode string) ([]string, error)
	Check(ctx context.Context, userID int64, appCode string, permission string) (*dto.AccessCheckResponse, error)
	AccessToken(ctx context.Context, userID int64, appCode string) (*dto.AccessTokenResponse, error)
}
