package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/teamroles/dto"
	"github.com/open-suite/authorization/internal/modules/teamroles/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
)

type TeamRoleServiceImpl struct {
	TeamRoleRepository repositories.TeamRoleRepository
	log                *logger.LayerLogger
}

func NewTeamRoleService(repository repositories.TeamRoleRepository, appLogger *logger.Logger) TeamRoleService {
	return &TeamRoleServiceImpl{
		TeamRoleRepository: repository,
		log:                appLogger.Layer("service.team_roles"),
	}
}

func (s *TeamRoleServiceImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.TeamRole, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.TeamRoleRepository.Find(ctx, params)
	end(err, "count", len(items))
	return items, err
}

func (s *TeamRoleServiceImpl) FindByID(ctx context.Context, id int64) (*entities.TeamRole, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.TeamRoleRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *TeamRoleServiceImpl) Create(ctx context.Context, request dto.CreateTeamRoleRequest) (*entities.TeamRole, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.TeamRoleRepository.Create(ctx, entities.TeamRole{
		TeamId: request.TeamId,
		RoleId: request.RoleId,
		AppId:  request.AppId,
	})
	end(err)
	return item, err
}

func (s *TeamRoleServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateTeamRoleRequest) (*entities.TeamRole, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.TeamId != nil {
		data["team_id"] = *request.TeamId
	}
	if request.RoleId != nil {
		data["role_id"] = *request.RoleId
	}
	if request.AppId != nil {
		data["app_id"] = *request.AppId
	}

	item, err := s.TeamRoleRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *TeamRoleServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.TeamRoleRepository.Delete(ctx, id)
	end(err)
	return err
}
