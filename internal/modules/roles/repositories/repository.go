package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type RoleRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.Role, error)
	FindByID(ctx context.Context, id int64) (*entities.Role, error)
	Create(ctx context.Context, entity entities.Role) (*entities.Role, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.Role, error)
	Delete(ctx context.Context, id int64) error
}
