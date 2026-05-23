package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type MenuRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.Menu, error)
	FindByID(ctx context.Context, id int64) (*entities.Menu, error)
	FindByCode(ctx context.Context, code string) (*entities.Menu, error)
	FindByAppID(ctx context.Context, appID int64, params shared.ListParams) ([]entities.Menu, error)
	FindByAppCode(ctx context.Context, appCode string, params shared.ListParams) ([]entities.Menu, error)
	Create(ctx context.Context, entity entities.Menu) (*entities.Menu, error)
	CreateBulk(ctx context.Context, items []entities.Menu) ([]entities.Menu, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.Menu, error)
	Delete(ctx context.Context, id int64) error
}
