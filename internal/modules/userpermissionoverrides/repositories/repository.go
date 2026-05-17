package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
)

type UserPermissionOverrideRepository interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.UserPermissionOverride, error)
	FindByID(ctx context.Context, id int64) (*entities.UserPermissionOverride, error)
	Create(ctx context.Context, entity entities.UserPermissionOverride) (*entities.UserPermissionOverride, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.UserPermissionOverride, error)
	Delete(ctx context.Context, id int64) error
}
