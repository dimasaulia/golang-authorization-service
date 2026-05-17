package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/teammembers/dto"
)

type TeamMemberService interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.TeamMember, error)
	FindByID(ctx context.Context, id int64) (*entities.TeamMember, error)
	Create(ctx context.Context, request dto.CreateTeamMemberRequest) (*entities.TeamMember, error)
	Update(ctx context.Context, id int64, request dto.UpdateTeamMemberRequest) (*entities.TeamMember, error)
	Delete(ctx context.Context, id int64) error
}
