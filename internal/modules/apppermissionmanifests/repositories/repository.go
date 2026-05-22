package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type AppPermissionManifestRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.AppPermissionManifest, error)
	FindByID(ctx context.Context, id int64) (*entities.AppPermissionManifest, error)
	Create(ctx context.Context, entity entities.AppPermissionManifest) (*entities.AppPermissionManifest, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.AppPermissionManifest, error)
	Delete(ctx context.Context, id int64) error
}
