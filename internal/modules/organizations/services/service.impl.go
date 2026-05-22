package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/organizations/dto"
	"github.com/open-suite/authorization/internal/modules/organizations/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
)

type OrganizationServiceImpl struct {
	OrganizationRepository repositories.OrganizationRepository
	log                    *logger.LayerLogger
}

func NewOrganizationService(repository repositories.OrganizationRepository, appLogger *logger.Logger) OrganizationService {
	return &OrganizationServiceImpl{
		OrganizationRepository: repository,
		log:                    appLogger.Layer("service.organizations"),
	}
}

func (s *OrganizationServiceImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.Organization, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.OrganizationRepository.Find(ctx, params)
	end(err, "count", len(items))
	return items, err
}

func (s *OrganizationServiceImpl) FindByID(ctx context.Context, id int64) (*entities.Organization, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.OrganizationRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *OrganizationServiceImpl) Create(ctx context.Context, request dto.CreateOrganizationRequest) (*entities.Organization, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.OrganizationRepository.Create(ctx, entities.Organization{
		Code:   request.Code,
		Name:   request.Name,
		Type:   request.Type,
		Status: request.Status,
	})
	end(err)
	return item, err
}

func (s *OrganizationServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateOrganizationRequest) (*entities.Organization, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.Code != nil {
		data["code"] = *request.Code
	}
	if request.Name != nil {
		data["name"] = *request.Name
	}
	if request.Type != nil {
		data["type"] = *request.Type
	}
	if request.Status != nil {
		data["status"] = *request.Status
	}

	item, err := s.OrganizationRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *OrganizationServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.OrganizationRepository.Delete(ctx, id)
	end(err)
	return err
}
