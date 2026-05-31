package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/teammembers/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type TeamMemberService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.TeamMember, error)
	FindByID(ctx context.Context, id int64) (*entities.TeamMember, error)
	Create(ctx context.Context, request dto.CreateTeamMemberRequest) (*entities.TeamMember, error)
	AssignTeamsToUser(ctx context.Context, userID int64, teamIDs []int64) ([]entities.TeamMember, error)
	FindTeamIDsByUserID(ctx context.Context, userID int64) ([]int64, error)
	FindAssignedTeamsByUserID(ctx context.Context, userID int64) ([]entities.UserAssignedTeam, error)
	ReplaceTeamsForUser(ctx context.Context, userID int64, teamIDs []int64) ([]entities.TeamMember, error)
	Update(ctx context.Context, id int64, request dto.UpdateTeamMemberRequest) (*entities.TeamMember, error)
	Delete(ctx context.Context, id int64) error
}
