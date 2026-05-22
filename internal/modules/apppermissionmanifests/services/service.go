package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/apppermissionmanifests/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type AppPermissionManifestService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.AppPermissionManifest, error)
	FindByID(ctx context.Context, id int64) (*entities.AppPermissionManifest, error)
	Create(ctx context.Context, request dto.CreateAppPermissionManifestRequest) (*entities.AppPermissionManifest, error)
	Update(ctx context.Context, id int64, request dto.UpdateAppPermissionManifestRequest) (*entities.AppPermissionManifest, error)
	Delete(ctx context.Context, id int64) error
}
