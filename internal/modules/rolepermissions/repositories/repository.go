package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type RolePermissionRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.RolePermission, error)
	FindByAppID(ctx context.Context, appID int64, params shared.ListParams) ([]entities.RolePermissionDetail, error)
	FindByAppCode(ctx context.Context, appCode string, params shared.ListParams) ([]entities.RolePermissionDetail, error)
	FindByRole(ctx context.Context, roleID int64, params shared.ListParams) ([]entities.RolePermissionDetail, error)
	FindRoleIDByCode(ctx context.Context, roleCode string) (int64, error)
	FindUserIDsByRoleID(ctx context.Context, roleID int64) ([]int64, error)
	FindAppCodesByIDs(ctx context.Context, appIDs []int64) ([]string, error)
	FindRoleSummaries(ctx context.Context, params shared.ListParams) ([]entities.RolePermissionSummary, error)
	FindRoleSummariesByAppID(ctx context.Context, appID int64, params shared.ListParams) ([]entities.RolePermissionSummary, error)
	FindRoleSummariesByAppCode(ctx context.Context, appCode string, params shared.ListParams) ([]entities.RolePermissionSummary, error)
	FindAvailablePermissionsByAppID(ctx context.Context, appID int64, params shared.ListParams) ([]entities.AvailablePermissionRow, error)
	FindAvailablePermissionsByAppCode(ctx context.Context, appCode string, params shared.ListParams) ([]entities.AvailablePermissionRow, error)
	FindByID(ctx context.Context, id int64) (*entities.RolePermission, error)
	Create(ctx context.Context, entity entities.RolePermission) (*entities.RolePermission, error)
	CreateBulk(ctx context.Context, items []entities.RolePermission) ([]entities.RolePermission, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.RolePermission, error)
	ReplaceByRoleAndApps(ctx context.Context, roleID int64, appIDs []int64, items []entities.RolePermission) ([]entities.RolePermission, error)
	Delete(ctx context.Context, id int64) error
}
