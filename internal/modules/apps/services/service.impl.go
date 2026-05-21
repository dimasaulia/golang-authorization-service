package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/apps/dto"
	"github.com/open-suite/authorization/internal/modules/apps/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
)

type AppServiceImpl struct {
	AppRepository repositories.AppRepository
	log           *logger.LayerLogger
}

func NewAppService(repository repositories.AppRepository, appLogger *logger.Logger) AppService {
	return &AppServiceImpl{
		AppRepository: repository,
		log:           appLogger.Layer("service.apps"),
	}
}

func (s *AppServiceImpl) Find(ctx context.Context, limit uint64, offset uint64, search string) ([]entities.App, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.AppRepository.Find(ctx, limit, offset, search)
	end(err, "count", len(items))
	return items, err
}

func (s *AppServiceImpl) FindByID(ctx context.Context, id int64) (*entities.App, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.AppRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *AppServiceImpl) Create(ctx context.Context, request dto.CreateAppRequest) (*entities.App, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.AppRepository.Create(ctx, entities.App{
		Code:   request.Code,
		Name:   request.Name,
		Status: request.Status,
	})
	end(err)
	return item, err
}

func (s *AppServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateAppRequest) (*entities.App, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.Code != nil {
		data["code"] = *request.Code
	}
	if request.Name != nil {
		data["name"] = *request.Name
	}
	if request.Status != nil {
		data["status"] = *request.Status
	}

	item, err := s.AppRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *AppServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.AppRepository.Delete(ctx, id)
	end(err)
	return err
}
