package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type AccessCacheVersionRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.AccessCacheVersion, error)
	FindByID(ctx context.Context, id int64) (*entities.AccessCacheVersion, error)
	Create(ctx context.Context, entity entities.AccessCacheVersion) (*entities.AccessCacheVersion, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.AccessCacheVersion, error)
	Delete(ctx context.Context, id int64) error
}
