package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type ModuleRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.Module, error)
	FindByID(ctx context.Context, id int64) (*entities.Module, error)
	FindByCode(ctx context.Context, code string) (*entities.Module, error)
	FindByAppID(ctx context.Context, appID int64, params shared.ListParams) ([]entities.Module, error)
	FindByAppCode(ctx context.Context, appCode string, params shared.ListParams) ([]entities.Module, error)
	Create(ctx context.Context, entity entities.Module) (*entities.Module, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.Module, error)
	Delete(ctx context.Context, id int64) error
}
