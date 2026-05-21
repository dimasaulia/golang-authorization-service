package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
)

type ActionRepository interface {
	Find(ctx context.Context, limit uint64, offset uint64, search string) ([]entities.Action, error)
	FindByID(ctx context.Context, id int64) (*entities.Action, error)
	Create(ctx context.Context, entity entities.Action) (*entities.Action, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.Action, error)
	Delete(ctx context.Context, id int64) error
}
