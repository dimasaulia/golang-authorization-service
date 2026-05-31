package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/userroles/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type UserRoleService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.UserRole, error)
	FindByID(ctx context.Context, id int64) (*entities.UserRole, error)
	Create(ctx context.Context, request dto.CreateUserRoleRequest) (*entities.UserRole, error)
	AssignRolesToUser(ctx context.Context, userID int64, roleIDs []int64, organizationID *int64, assignedBy *int64) ([]entities.UserRole, error)
	FindRoleIDsByUserID(ctx context.Context, userID int64) ([]int64, error)
	FindAssignedRolesByUserID(ctx context.Context, userID int64) ([]entities.UserAssignedRole, error)
	ReplaceRolesForUser(ctx context.Context, userID int64, roleIDs []int64, organizationID *int64, assignedBy *int64) ([]entities.UserRole, error)
	Update(ctx context.Context, id int64, request dto.UpdateUserRoleRequest) (*entities.UserRole, error)
	Delete(ctx context.Context, id int64) error
}
