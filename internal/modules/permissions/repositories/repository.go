package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
)

type PermissionRepository interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.Permission, error)
	FindByID(ctx context.Context, id int64) (*entities.Permission, error)
	FindByCode(ctx context.Context, code string) (*entities.Permission, error)
	Create(ctx context.Context, entity entities.Permission) (*entities.Permission, error)
	CreateBulk(ctx context.Context, items []entities.Permission) ([]entities.Permission, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.Permission, error)
	Delete(ctx context.Context, id int64) error
}
