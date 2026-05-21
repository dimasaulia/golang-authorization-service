package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/teams/dto"
)

type TeamService interface {
	Find(ctx context.Context, limit uint64, offset uint64, search string) ([]entities.Team, error)
	FindByID(ctx context.Context, id int64) (*entities.Team, error)
	Create(ctx context.Context, request dto.CreateTeamRequest) (*entities.Team, error)
	Update(ctx context.Context, id int64, request dto.UpdateTeamRequest) (*entities.Team, error)
	Delete(ctx context.Context, id int64) error
}
