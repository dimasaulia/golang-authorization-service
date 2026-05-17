package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/permissions/dto"
)

type PermissionService interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.Permission, error)
	FindByID(ctx context.Context, id int64) (*entities.Permission, error)
	Create(ctx context.Context, request dto.CreatePermissionRequest) (*entities.Permission, error)
	Update(ctx context.Context, id int64, request dto.UpdatePermissionRequest) (*entities.Permission, error)
	Delete(ctx context.Context, id int64) error
}
