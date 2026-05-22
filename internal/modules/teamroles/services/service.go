package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/teamroles/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type TeamRoleService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.TeamRole, error)
	FindByID(ctx context.Context, id int64) (*entities.TeamRole, error)
	Create(ctx context.Context, request dto.CreateTeamRoleRequest) (*entities.TeamRole, error)
	Update(ctx context.Context, id int64, request dto.UpdateTeamRoleRequest) (*entities.TeamRole, error)
	Delete(ctx context.Context, id int64) error
}
