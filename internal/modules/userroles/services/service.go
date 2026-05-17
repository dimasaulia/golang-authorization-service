package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/userroles/dto"
)

type UserRoleService interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.UserRole, error)
	FindByID(ctx context.Context, id int64) (*entities.UserRole, error)
	Create(ctx context.Context, request dto.CreateUserRoleRequest) (*entities.UserRole, error)
	Update(ctx context.Context, id int64, request dto.UpdateUserRoleRequest) (*entities.UserRole, error)
	Delete(ctx context.Context, id int64) error
}
