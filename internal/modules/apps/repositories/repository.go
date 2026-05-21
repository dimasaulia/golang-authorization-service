package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
)

type AppRepository interface {
	Find(ctx context.Context, limit uint64, offset uint64, search string) ([]entities.App, error)
	FindByID(ctx context.Context, id int64) (*entities.App, error)
	Create(ctx context.Context, entity entities.App) (*entities.App, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.App, error)
	Delete(ctx context.Context, id int64) error
}
