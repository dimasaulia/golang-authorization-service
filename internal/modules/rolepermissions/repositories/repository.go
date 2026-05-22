package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type RolePermissionRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.RolePermission, error)
	FindByID(ctx context.Context, id int64) (*entities.RolePermission, error)
	Create(ctx context.Context, entity entities.RolePermission) (*entities.RolePermission, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.RolePermission, error)
	Delete(ctx context.Context, id int64) error
}
