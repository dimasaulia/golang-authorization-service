package services

import (
	"context"
	"strconv"
	"strings"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/roles/dto"
	"github.com/open-suite/authorization/internal/modules/roles/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
)

type RoleServiceImpl struct {
	RoleRepository repositories.RoleRepository
	log            *logger.LayerLogger
}

func NewRoleService(repository repositories.RoleRepository, appLogger *logger.Logger) RoleService {
	return &RoleServiceImpl{
		RoleRepository: repository,
		log:            appLogger.Layer("service.roles"),
	}
}

func (s *RoleServiceImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.Role, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.RoleRepository.Find(ctx, params)
	end(err, "count", len(items))
	return items, err
}

func (s *RoleServiceImpl) FindByID(ctx context.Context, id int64) (*entities.Role, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.RoleRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *RoleServiceImpl) FindByUnique(ctx context.Context, unique string) (*entities.Role, error) {
	unique = strings.TrimSpace(unique)
	end := s.log.Start(ctx, "FindByUnique", "unique", unique)

	id, err := strconv.ParseInt(unique, 10, 64)
	if err == nil {
		item, err := s.RoleRepository.FindByID(ctx, id)
		end(err)
		return item, err
	}

	item, err := s.RoleRepository.FindByCode(ctx, unique)
	end(err)
	return item, err
}

func (s *RoleServiceImpl) FindByCode(ctx context.Context, code string) (*entities.Role, error) {
	end := s.log.Start(ctx, "FindByCode", "code", code)
	item, err := s.RoleRepository.FindByCode(ctx, code)
	end(err)
	return item, err
}

func (s *RoleServiceImpl) FindByApp(ctx context.Context, appUnique string, params shared.ListParams) ([]entities.Role, error) {
	appUnique = strings.TrimSpace(appUnique)
	end := s.log.Start(ctx, "FindByApp", "app_unique", appUnique)

	id, err := strconv.ParseInt(appUnique, 10, 64)
	if err == nil {
		items, err := s.RoleRepository.FindByAppID(ctx, id, params)
		end(err, "count", len(items))
		return items, err
	}

	items, err := s.RoleRepository.FindByAppCode(ctx, appUnique, params)
	end(err, "count", len(items))
	return items, err
}

func (s *RoleServiceImpl) Create(ctx context.Context, request dto.CreateRoleRequest) (*entities.Role, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.RoleRepository.Create(ctx, entities.Role{
		OrganizationId: request.OrganizationId,
		AppId:          request.AppId,
		Code:           request.Code,
		Name:           request.Name,
		Description:    request.Description,
		Scope:          request.Scope,
		IsSystem:       request.IsSystem,
		Status:         request.Status,
	})
	end(err)
	return item, err
}

func (s *RoleServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateRoleRequest) (*entities.Role, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.OrganizationId != nil {
		data["organization_id"] = *request.OrganizationId
	}
	if request.AppId != nil {
		data["app_id"] = *request.AppId
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
	if request.Scope != nil {
		data["scope"] = *request.Scope
	}
	if request.IsSystem != nil {
		data["is_system"] = *request.IsSystem
	}
	if request.Status != nil {
		data["status"] = *request.Status
	}

	item, err := s.RoleRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *RoleServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.RoleRepository.Delete(ctx, id)
	end(err)
	return err
}
