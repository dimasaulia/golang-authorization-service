package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
)

type AppClientRepository interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.AppClient, error)
	FindByID(ctx context.Context, id int64) (*entities.AppClient, error)
	Create(ctx context.Context, entity entities.AppClient) (*entities.AppClient, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.AppClient, error)
	Delete(ctx context.Context, id int64) error
}
