package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type UserRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.User, error)
	FindByID(ctx context.Context, id int64) (*entities.User, error)
	Create(ctx context.Context, entity entities.User) (*entities.User, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.User, error)
	Delete(ctx context.Context, id int64) error
}
