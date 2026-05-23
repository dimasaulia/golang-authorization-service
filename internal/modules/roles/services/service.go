package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/roles/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type RoleService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.Role, error)
	FindByID(ctx context.Context, id int64) (*entities.Role, error)
	FindByUnique(ctx context.Context, unique string) (*entities.Role, error)
	FindByCode(ctx context.Context, code string) (*entities.Role, error)
	FindByApp(ctx context.Context, appUnique string, params shared.ListParams) ([]entities.Role, error)
	Create(ctx context.Context, request dto.CreateRoleRequest) (*entities.Role, error)
	Update(ctx context.Context, id int64, request dto.UpdateRoleRequest) (*entities.Role, error)
	Delete(ctx context.Context, id int64) error
}
