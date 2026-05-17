package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
)

type ModuleRepository interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.Module, error)
	FindByID(ctx context.Context, id int64) (*entities.Module, error)
	Create(ctx context.Context, entity entities.Module) (*entities.Module, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.Module, error)
	Delete(ctx context.Context, id int64) error
}
