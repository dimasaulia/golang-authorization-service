package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
)

type RoleRepository interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.Role, error)
	FindByID(ctx context.Context, id int64) (*entities.Role, error)
	Create(ctx context.Context, entity entities.Role) (*entities.Role, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.Role, error)
	Delete(ctx context.Context, id int64) error
}
