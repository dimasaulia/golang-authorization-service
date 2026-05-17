package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/teammembers/dto"
	"github.com/open-suite/authorization/internal/modules/teammembers/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
)

type TeamMemberServiceImpl struct {
	TeamMemberRepository repositories.TeamMemberRepository
	log                  *logger.LayerLogger
}

func NewTeamMemberService(repository repositories.TeamMemberRepository, appLogger *logger.Logger) TeamMemberService {
	return &TeamMemberServiceImpl{
		TeamMemberRepository: repository,
		log:                  appLogger.Layer("service.team_members"),
	}
}

func (s *TeamMemberServiceImpl) Find(ctx context.Context, limit uint64, offset uint64) ([]entities.TeamMember, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.TeamMemberRepository.Find(ctx, limit, offset)
	end(err, "count", len(items))
	return items, err
}

func (s *TeamMemberServiceImpl) FindByID(ctx context.Context, id int64) (*entities.TeamMember, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.TeamMemberRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *TeamMemberServiceImpl) Create(ctx context.Context, request dto.CreateTeamMemberRequest) (*entities.TeamMember, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.TeamMemberRepository.Create(ctx, entities.TeamMember{
		TeamId:     request.TeamId,
		UserId:     request.UserId,
		RoleInTeam: request.RoleInTeam,
	})
	end(err)
	return item, err
}

func (s *TeamMemberServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateTeamMemberRequest) (*entities.TeamMember, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.TeamId != nil {
		data["team_id"] = *request.TeamId
	}
	if request.UserId != nil {
		data["user_id"] = *request.UserId
	}
	if request.RoleInTeam != nil {
		data["role_in_team"] = *request.RoleInTeam
	}

	item, err := s.TeamMemberRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *TeamMemberServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.TeamMemberRepository.Delete(ctx, id)
	end(err)
	return err
}
