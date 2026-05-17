package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/apppermissionmanifests/dto"
)

type AppPermissionManifestService interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.AppPermissionManifest, error)
	FindByID(ctx context.Context, id int64) (*entities.AppPermissionManifest, error)
	Create(ctx context.Context, request dto.CreateAppPermissionManifestRequest) (*entities.AppPermissionManifest, error)
	Update(ctx context.Context, id int64, request dto.UpdateAppPermissionManifestRequest) (*entities.AppPermissionManifest, error)
	Delete(ctx context.Context, id int64) error
}
