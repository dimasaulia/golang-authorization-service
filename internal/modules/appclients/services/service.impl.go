package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/appclients/dto"
	"github.com/open-suite/authorization/internal/modules/appclients/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
)

type AppClientServiceImpl struct {
	AppClientRepository repositories.AppClientRepository
	log                 *logger.LayerLogger
}

func NewAppClientService(repository repositories.AppClientRepository, appLogger *logger.Logger) AppClientService {
	return &AppClientServiceImpl{
		AppClientRepository: repository,
		log:                 appLogger.Layer("service.app_clients"),
	}
}

func (s *AppClientServiceImpl) Find(ctx context.Context, limit uint64, offset uint64) ([]entities.AppClient, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.AppClientRepository.Find(ctx, limit, offset)
	end(err, "count", len(items))
	return items, err
}

func (s *AppClientServiceImpl) FindByID(ctx context.Context, id int64) (*entities.AppClient, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.AppClientRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *AppClientServiceImpl) Create(ctx context.Context, request dto.CreateAppClientRequest) (*entities.AppClient, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.AppClientRepository.Create(ctx, entities.AppClient{
		AppId:            request.AppId,
		KeycloakClientId: request.KeycloakClientId,
		Name:             request.Name,
		Status:           request.Status,
	})
	end(err)
	return item, err
}

func (s *AppClientServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateAppClientRequest) (*entities.AppClient, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.AppId != nil {
		data["app_id"] = *request.AppId
	}
	if request.KeycloakClientId != nil {
		data["keycloak_client_id"] = *request.KeycloakClientId
	}
	if request.Name != nil {
		data["name"] = *request.Name
	}
	if request.Status != nil {
		data["status"] = *request.Status
	}

	item, err := s.AppClientRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *AppClientServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.AppClientRepository.Delete(ctx, id)
	end(err)
	return err
}
