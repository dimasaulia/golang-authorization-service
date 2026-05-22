package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type ActionRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.Action, error)
	FindByID(ctx context.Context, id int64) (*entities.Action, error)
	FindByCode(ctx context.Context, code string) (*entities.Action, error)
	Create(ctx context.Context, entity entities.Action) (*entities.Action, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.Action, error)
	Delete(ctx context.Context, id int64) error
}
