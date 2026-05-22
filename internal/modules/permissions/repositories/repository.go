package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type PermissionRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.Permission, error)
	FindByID(ctx context.Context, id int64) (*entities.Permission, error)
	FindByCode(ctx context.Context, code string) (*entities.Permission, error)
	FindByAppID(ctx context.Context, appID int64, params shared.ListParams) ([]entities.Permission, error)
	FindByAppCode(ctx context.Context, appCode string, params shared.ListParams) ([]entities.Permission, error)
	Create(ctx context.Context, entity entities.Permission) (*entities.Permission, error)
	CreateBulk(ctx context.Context, items []entities.Permission) ([]entities.Permission, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.Permission, error)
	Delete(ctx context.Context, id int64) error
}
