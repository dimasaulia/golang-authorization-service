package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/rolepermissions/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type RolePermissionService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.RolePermission, error)
	FindByID(ctx context.Context, id int64) (*entities.RolePermission, error)
	Create(ctx context.Context, request dto.CreateRolePermissionRequest) (*entities.RolePermission, error)
	Update(ctx context.Context, id int64, request dto.UpdateRolePermissionRequest) (*entities.RolePermission, error)
	Delete(ctx context.Context, id int64) error
}
