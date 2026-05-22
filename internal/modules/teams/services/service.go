package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/teams/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type TeamService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.Team, error)
	FindByID(ctx context.Context, id int64) (*entities.Team, error)
	Create(ctx context.Context, request dto.CreateTeamRequest) (*entities.Team, error)
	Update(ctx context.Context, id int64, request dto.UpdateTeamRequest) (*entities.Team, error)
	Delete(ctx context.Context, id int64) error
}
