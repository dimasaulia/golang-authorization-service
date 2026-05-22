package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type AppClientRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.AppClient, error)
	FindByID(ctx context.Context, id int64) (*entities.AppClient, error)
	Create(ctx context.Context, entity entities.AppClient) (*entities.AppClient, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.AppClient, error)
	Delete(ctx context.Context, id int64) error
}
