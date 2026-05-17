package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
)

type UserIdentityRepository interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.UserIdentity, error)
	FindByID(ctx context.Context, id int64) (*entities.UserIdentity, error)
	Create(ctx context.Context, entity entities.UserIdentity) (*entities.UserIdentity, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.UserIdentity, error)
	Delete(ctx context.Context, id int64) error
}
