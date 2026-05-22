package services

import (
	"context"
	"strconv"
	"strings"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/teams/dto"
	"github.com/open-suite/authorization/internal/modules/teams/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
)

type TeamServiceImpl struct {
	TeamRepository repositories.TeamRepository
	log            *logger.LayerLogger
}

func NewTeamService(repository repositories.TeamRepository, appLogger *logger.Logger) TeamService {
	return &TeamServiceImpl{
		TeamRepository: repository,
		log:            appLogger.Layer("service.teams"),
	}
}

func (s *TeamServiceImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.Team, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.TeamRepository.Find(ctx, params)
	end(err, "count", len(items))
	return items, err
}

func (s *TeamServiceImpl) FindByID(ctx context.Context, id int64) (*entities.Team, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.TeamRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *TeamServiceImpl) FindByUnique(ctx context.Context, unique string) (*entities.Team, error) {
	unique = strings.TrimSpace(unique)
	end := s.log.Start(ctx, "FindByUnique", "unique", unique)

	id, err := strconv.ParseInt(unique, 10, 64)
	if err == nil {
		item, err := s.TeamRepository.FindByID(ctx, id)
		end(err)
		return item, err
	}

	item, err := s.TeamRepository.FindByCode(ctx, unique)
	end(err)
	return item, err
}

func (s *TeamServiceImpl) FindByCode(ctx context.Context, code string) (*entities.Team, error) {
	end := s.log.Start(ctx, "FindByCode", "code", code)
	item, err := s.TeamRepository.FindByCode(ctx, code)
	end(err)
	return item, err
}

func (s *TeamServiceImpl) Create(ctx context.Context, request dto.CreateTeamRequest) (*entities.Team, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.TeamRepository.Create(ctx, entities.Team{
		OrganizationId: request.OrganizationId,
		Code:           request.Code,
		Name:           request.Name,
		Status:         request.Status,
	})
	end(err)
	return item, err
}

func (s *TeamServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateTeamRequest) (*entities.Team, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.OrganizationId != nil {
		data["organization_id"] = *request.OrganizationId
	}
	if request.Code != nil {
		data["code"] = *request.Code
	}
	if request.Name != nil {
		data["name"] = *request.Name
	}
	if request.Status != nil {
		data["status"] = *request.Status
	}

	item, err := s.TeamRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *TeamServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.TeamRepository.Delete(ctx, id)
	end(err)
	return err
}
