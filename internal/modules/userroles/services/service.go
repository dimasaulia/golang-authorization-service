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
	Update(ctx context.Context, id int64, request dto.UpdateUserRoleRequest) (*entities.UserRole, error)
	Delete(ctx context.Context, id int64) error
}
