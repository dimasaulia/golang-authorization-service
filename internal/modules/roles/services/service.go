package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/roles/dto"
)

type RoleService interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.Role, error)
	FindByID(ctx context.Context, id int64) (*entities.Role, error)
	Create(ctx context.Context, request dto.CreateRoleRequest) (*entities.Role, error)
	Update(ctx context.Context, id int64, request dto.UpdateRoleRequest) (*entities.Role, error)
	Delete(ctx context.Context, id int64) error
}
