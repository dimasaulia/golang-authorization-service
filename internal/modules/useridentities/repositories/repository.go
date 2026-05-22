package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type UserIdentityRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.UserIdentity, error)
	FindByID(ctx context.Context, id int64) (*entities.UserIdentity, error)
	Create(ctx context.Context, entity entities.UserIdentity) (*entities.UserIdentity, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.UserIdentity, error)
	Delete(ctx context.Context, id int64) error
}
