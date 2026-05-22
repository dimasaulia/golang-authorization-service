package services

import (
	"context"
	"strconv"
	"strings"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/modules/dto"
	"github.com/open-suite/authorization/internal/modules/modules/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
)

type ModuleServiceImpl struct {
	ModuleRepository repositories.ModuleRepository
	log              *logger.LayerLogger
}

func NewModuleService(repository repositories.ModuleRepository, appLogger *logger.Logger) ModuleService {
	return &ModuleServiceImpl{
		ModuleRepository: repository,
		log:              appLogger.Layer("service.modules"),
	}
}

func (s *ModuleServiceImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.Module, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.ModuleRepository.Find(ctx, params)
	end(err, "count", len(items))
	return items, err
}

func (s *ModuleServiceImpl) FindByID(ctx context.Context, id int64) (*entities.Module, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.ModuleRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *ModuleServiceImpl) FindByUnique(ctx context.Context, unique string) (*entities.Module, error) {
	unique = strings.TrimSpace(unique)
	end := s.log.Start(ctx, "FindByUnique", "unique", unique)

	id, err := strconv.ParseInt(unique, 10, 64)
	if err == nil {
		item, err := s.ModuleRepository.FindByID(ctx, id)
		end(err)
		return item, err
	}

	item, err := s.ModuleRepository.FindByCode(ctx, unique)
	end(err)
	return item, err
}

func (s *ModuleServiceImpl) FindByCode(ctx context.Context, code string) (*entities.Module, error) {
	end := s.log.Start(ctx, "FindByCode", "code", code)
	item, err := s.ModuleRepository.FindByCode(ctx, code)
	end(err)
	return item, err
}

func (s *ModuleServiceImpl) FindByApp(ctx context.Context, appUnique string, params shared.ListParams) ([]entities.Module, error) {
	appUnique = strings.TrimSpace(appUnique)
	end := s.log.Start(ctx, "FindByApp", "app_unique", appUnique)

	id, err := strconv.ParseInt(appUnique, 10, 64)
	if err == nil {
		items, err := s.ModuleRepository.FindByAppID(ctx, id, params)
		end(err, "count", len(items))
		return items, err
	}

	items, err := s.ModuleRepository.FindByAppCode(ctx, appUnique, params)
	end(err, "count", len(items))
	return items, err
}

func (s *ModuleServiceImpl) Create(ctx context.Context, request dto.CreateModuleRequest) (*entities.Module, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.ModuleRepository.Create(ctx, entities.Module{
		AppId:  request.AppId,
		Code:   request.Code,
		Name:   request.Name,
		Status: request.Status,
	})
	end(err)
	return item, err
}

func (s *ModuleServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateModuleRequest) (*entities.Module, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.AppId != nil {
		data["app_id"] = *request.AppId
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

	item, err := s.ModuleRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *ModuleServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.ModuleRepository.Delete(ctx, id)
	end(err)
	return err
}
