package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type AppRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.App, error)
	FindByID(ctx context.Context, id int64) (*entities.App, error)
	FindByCode(ctx context.Context, code string) (*entities.App, error)
	Create(ctx context.Context, entity entities.App) (*entities.App, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.App, error)
	Delete(ctx context.Context, id int64) error
}
