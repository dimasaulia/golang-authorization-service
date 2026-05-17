package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/rolepermissions/dto"
)

type RolePermissionService interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.RolePermission, error)
	FindByID(ctx context.Context, id int64) (*entities.RolePermission, error)
	Create(ctx context.Context, request dto.CreateRolePermissionRequest) (*entities.RolePermission, error)
	Update(ctx context.Context, id int64, request dto.UpdateRolePermissionRequest) (*entities.RolePermission, error)
	Delete(ctx context.Context, id int64) error
}
