package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
)

type UserRepository interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.User, error)
	FindByID(ctx context.Context, id int64) (*entities.User, error)
	Create(ctx context.Context, entity entities.User) (*entities.User, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.User, error)
	Delete(ctx context.Context, id int64) error
}
