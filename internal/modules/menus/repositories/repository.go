package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type MenuRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.Menu, error)
	FindByID(ctx context.Context, id int64) (*entities.Menu, error)
	Create(ctx context.Context, entity entities.Menu) (*entities.Menu, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.Menu, error)
	Delete(ctx context.Context, id int64) error
}
