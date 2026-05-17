package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
)

type UserRoleRepository interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.UserRole, error)
	FindByID(ctx context.Context, id int64) (*entities.UserRole, error)
	Create(ctx context.Context, entity entities.UserRole) (*entities.UserRole, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.UserRole, error)
	Delete(ctx context.Context, id int64) error
}
