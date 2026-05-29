package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/modules/auth/dto"
)

type AccessRepository interface {
	FindAccessSummary(ctx context.Context, userID int64, appCode string) (*dto.AccessSummaryResponse, error)
	FindApps(ctx context.Context, userID int64) ([]dto.UserAppAccess, error)
	FindUserProviders(ctx context.Context, userID int64) ([]dto.CurrentUserProvider, error)
	FindAppPermissionCountSummaries(ctx context.Context, userID int64) ([]dto.AppPermissionCountSummary, error)
}
