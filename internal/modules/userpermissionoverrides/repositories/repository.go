package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type UserPermissionOverrideRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.UserPermissionOverride, error)
	FindByID(ctx context.Context, id int64) (*entities.UserPermissionOverride, error)
	Create(ctx context.Context, entity entities.UserPermissionOverride) (*entities.UserPermissionOverride, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.UserPermissionOverride, error)
	Delete(ctx context.Context, id int64) error
}
