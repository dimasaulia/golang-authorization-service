package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/rolepermissions/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type RolePermissionService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.RolePermission, error)
	FindByApp(ctx context.Context, appUnique string, params shared.ListParams) ([]entities.RolePermissionDetail, error)
	FindByRole(ctx context.Context, roleUnique string, params shared.ListParams) ([]entities.RolePermissionDetail, error)
	FindRoleSummaries(ctx context.Context, params shared.ListParams) ([]entities.RolePermissionSummary, error)
	FindRoleSummariesByApp(ctx context.Context, appUnique string, params shared.ListParams) ([]entities.RolePermissionSummary, error)
	FindAvailablePermissionsByApp(ctx context.Context, appUnique string, params shared.ListParams) ([]entities.AvailablePermissionModule, error)
	FindByID(ctx context.Context, id int64) (*entities.RolePermission, error)
	Create(ctx context.Context, request dto.CreateRolePermissionRequest) (*entities.RolePermission, error)
	CreateBulk(ctx context.Context, request []dto.CreateBulkRolePermissionRequest) ([]entities.RolePermission, error)
	Update(ctx context.Context, id int64, request dto.UpdateRolePermissionRequest) (*entities.RolePermission, error)
	UpdateByRole(ctx context.Context, roleUnique string, request dto.UpdateRolePermissionByRoleRequest) ([]entities.RolePermission, error)
	Delete(ctx context.Context, id int64) error
}
