package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/teammembers/dto"
	"github.com/open-suite/authorization/internal/modules/teammembers/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
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

func (s *TeamMemberServiceImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.TeamMember, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.TeamMemberRepository.Find(ctx, params)
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

func (s *TeamMemberServiceImpl) AssignTeamsToUser(ctx context.Context, userID int64, teamIDs []int64) ([]entities.TeamMember, error) {
	end := s.log.Start(ctx, "AssignTeamsToUser", "user_id", userID, "count", len(teamIDs))
	items := make([]entities.TeamMember, 0, len(teamIDs))
	for _, teamID := range teamIDs {
		if teamID == 0 {
			continue
		}
		item, err := s.TeamMemberRepository.Create(ctx, entities.TeamMember{
			TeamId: teamID,
			UserId: userID,
		})
		if err != nil {
			end(err)
			return nil, err
		}
		items = append(items, *item)
	}
	end(nil, "assigned", len(items))
	return items, nil
}

func (s *TeamMemberServiceImpl) FindTeamIDsByUserID(ctx context.Context, userID int64) ([]int64, error) {
	end := s.log.Start(ctx, "FindTeamIDsByUserID", "user_id", userID)
	items, err := s.TeamMemberRepository.FindTeamIDsByUserID(ctx, userID)
	end(err, "count", len(items))
	return items, err
}

func (s *TeamMemberServiceImpl) FindAssignedTeamsByUserID(ctx context.Context, userID int64) ([]entities.UserAssignedTeam, error) {
	end := s.log.Start(ctx, "FindAssignedTeamsByUserID", "user_id", userID)
	items, err := s.TeamMemberRepository.FindAssignedTeamsByUserID(ctx, userID)
	end(err, "count", len(items))
	return items, err
}

func (s *TeamMemberServiceImpl) ReplaceTeamsForUser(ctx context.Context, userID int64, teamIDs []int64) ([]entities.TeamMember, error) {
	end := s.log.Start(ctx, "ReplaceTeamsForUser", "user_id", userID, "count", len(teamIDs))
	items, err := s.TeamMemberRepository.ReplaceTeamsForUser(ctx, userID, teamIDs)
	end(err, "assigned", len(items))
	return items, err
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
