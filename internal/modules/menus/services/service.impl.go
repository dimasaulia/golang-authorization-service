package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/menus/dto"
	"github.com/open-suite/authorization/internal/modules/menus/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
)

type MenuServiceImpl struct {
	MenuRepository repositories.MenuRepository
	log            *logger.LayerLogger
}

func NewMenuService(repository repositories.MenuRepository, appLogger *logger.Logger) MenuService {
	return &MenuServiceImpl{
		MenuRepository: repository,
		log:            appLogger.Layer("service.menus"),
	}
}

func (s *MenuServiceImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.Menu, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.MenuRepository.Find(ctx, params)
	end(err, "count", len(items))
	return items, err
}

func (s *MenuServiceImpl) FindByID(ctx context.Context, id int64) (*entities.Menu, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.MenuRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *MenuServiceImpl) Create(ctx context.Context, request dto.CreateMenuRequest) (*entities.Menu, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.MenuRepository.Create(ctx, entities.Menu{
		AppId:                request.AppId,
		ModuleId:             request.ModuleId,
		ParentId:             request.ParentId,
		Code:                 request.Code,
		Name:                 request.Name,
		RoutePath:            request.RoutePath,
		SortOrder:            request.SortOrder,
		RequiredPermissionId: request.RequiredPermissionId,
		Status:               request.Status,
	})
	end(err)
	return item, err
}

func (s *MenuServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateMenuRequest) (*entities.Menu, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.AppId != nil {
		data["app_id"] = *request.AppId
	}
	if request.ModuleId != nil {
		data["module_id"] = *request.ModuleId
	}
	if request.ParentId != nil {
		data["parent_id"] = *request.ParentId
	}
	if request.Code != nil {
		data["code"] = *request.Code
	}
	if request.Name != nil {
		data["name"] = *request.Name
	}
	if request.RoutePath != nil {
		data["route_path"] = *request.RoutePath
	}
	if request.SortOrder != nil {
		data["sort_order"] = *request.SortOrder
	}
	if request.RequiredPermissionId != nil {
		data["required_permission_id"] = *request.RequiredPermissionId
	}
	if request.Status != nil {
		data["status"] = *request.Status
	}

	item, err := s.MenuRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *MenuServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.MenuRepository.Delete(ctx, id)
	end(err)
	return err
}
