package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/accesscacheversions/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type AccessCacheVersionService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.AccessCacheVersion, error)
	FindByID(ctx context.Context, id int64) (*entities.AccessCacheVersion, error)
	Create(ctx context.Context, request dto.CreateAccessCacheVersionRequest) (*entities.AccessCacheVersion, error)
	Update(ctx context.Context, id int64, request dto.UpdateAccessCacheVersionRequest) (*entities.AccessCacheVersion, error)
	Delete(ctx context.Context, id int64) error
}
