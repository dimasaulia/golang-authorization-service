package services

import (
	"context"
	"strconv"
	"strings"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/permissions/dto"
	"github.com/open-suite/authorization/internal/modules/permissions/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
)

type PermissionServiceImpl struct {
	PermissionRepository repositories.PermissionRepository
	log                  *logger.LayerLogger
}

func NewPermissionService(repository repositories.PermissionRepository, appLogger *logger.Logger) PermissionService {
	return &PermissionServiceImpl{
		PermissionRepository: repository,
		log:                  appLogger.Layer("service.permissions"),
	}
}

func (s *PermissionServiceImpl) Find(ctx context.Context, limit uint64, offset uint64) ([]entities.Permission, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.PermissionRepository.Find(ctx, limit, offset)
	end(err, "count", len(items))
	return items, err
}

func (s *PermissionServiceImpl) FindByID(ctx context.Context, id int64) (*entities.Permission, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.PermissionRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *PermissionServiceImpl) FindByUnique(ctx context.Context, unique string) (*entities.Permission, error) {
	unique = strings.TrimSpace(unique)
	end := s.log.Start(ctx, "FindByUnique", "unique", unique)

	id, err := strconv.ParseInt(unique, 10, 64)
	if err == nil {
		item, err := s.PermissionRepository.FindByID(ctx, id)
		end(err)
		return item, err
	}

	item, err := s.PermissionRepository.FindByCode(ctx, unique)
	end(err)
	return item, err
}

func (s *PermissionServiceImpl) FindByCode(ctx context.Context, code string) (*entities.Permission, error) {
	end := s.log.Start(ctx, "FindByCode", "code", code)
	item, err := s.PermissionRepository.FindByCode(ctx, code)
	end(err)
	return item, err
}

func (s *PermissionServiceImpl) Create(ctx context.Context, request dto.CreatePermissionRequest) (*entities.Permission, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.PermissionRepository.Create(ctx, entities.Permission{
		AppId:       request.AppId,
		ModuleId:    request.ModuleId,
		ActionId:    request.ActionId,
		Code:        request.Code,
		Name:        request.Name,
		Description: request.Description,
		RiskLevel:   request.RiskLevel,
		IsSystem:    request.IsSystem,
		Status:      request.Status,
	})
	end(err)
	return item, err
}

func (s *PermissionServiceImpl) CreateBulk(ctx context.Context, request []dto.CreatePermissionRequest) ([]entities.Permission, error) {
	end := s.log.Start(ctx, "CreateBulk", "count", len(request))

	items := make([]entities.Permission, 0, len(request))
	for _, item := range request {
		items = append(items, entities.Permission{
			AppId:       item.AppId,
			ModuleId:    item.ModuleId,
			ActionId:    item.ActionId,
			Code:        item.Code,
			Name:        item.Name,
			Description: item.Description,
			RiskLevel:   item.RiskLevel,
			IsSystem:    item.IsSystem,
			Status:      item.Status,
		})
	}

	created, err := s.PermissionRepository.CreateBulk(ctx, items)
	end(err, "count", len(created))
	return created, err
}

func (s *PermissionServiceImpl) Update(ctx context.Context, id int64, request dto.UpdatePermissionRequest) (*entities.Permission, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.AppId != nil {
		data["app_id"] = *request.AppId
	}
	if request.ModuleId != nil {
		data["module_id"] = *request.ModuleId
	}
	if request.ActionId != nil {
		data["action_id"] = *request.ActionId
	}
	if request.Code != nil {
		data["code"] = *request.Code
	}
	if request.Name != nil {
		data["name"] = *request.Name
	}
	if request.Description != nil {
		data["description"] = *request.Description
	}
	if request.RiskLevel != nil {
		data["risk_level"] = *request.RiskLevel
	}
	if request.IsSystem != nil {
		data["is_system"] = *request.IsSystem
	}
	if request.Status != nil {
		data["status"] = *request.Status
	}

	item, err := s.PermissionRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *PermissionServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.PermissionRepository.Delete(ctx, id)
	end(err)
	return err
}
