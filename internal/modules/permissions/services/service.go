package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/permissions/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type PermissionService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.Permission, error)
	FindByID(ctx context.Context, id int64) (*entities.Permission, error)
	FindByUnique(ctx context.Context, unique string) (*entities.Permission, error)
	FindByCode(ctx context.Context, code string) (*entities.Permission, error)
	FindByApp(ctx context.Context, appUnique string, params shared.ListParams) ([]entities.Permission, error)
	Create(ctx context.Context, request dto.CreatePermissionRequest) (*entities.Permission, error)
	CreateBulk(ctx context.Context, request []dto.CreatePermissionRequest) ([]entities.Permission, error)
	Update(ctx context.Context, id int64, request dto.UpdatePermissionRequest) (*entities.Permission, error)
	Delete(ctx context.Context, id int64) error
}
