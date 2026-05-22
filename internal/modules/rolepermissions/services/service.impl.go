package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/rolepermissions/dto"
	"github.com/open-suite/authorization/internal/modules/rolepermissions/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
)

type RolePermissionServiceImpl struct {
	RolePermissionRepository repositories.RolePermissionRepository
	log                      *logger.LayerLogger
}

func NewRolePermissionService(repository repositories.RolePermissionRepository, appLogger *logger.Logger) RolePermissionService {
	return &RolePermissionServiceImpl{
		RolePermissionRepository: repository,
		log:                      appLogger.Layer("service.role_permissions"),
	}
}

func (s *RolePermissionServiceImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.RolePermission, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.RolePermissionRepository.Find(ctx, params)
	end(err, "count", len(items))
	return items, err
}

func (s *RolePermissionServiceImpl) FindByID(ctx context.Context, id int64) (*entities.RolePermission, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.RolePermissionRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *RolePermissionServiceImpl) Create(ctx context.Context, request dto.CreateRolePermissionRequest) (*entities.RolePermission, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.RolePermissionRepository.Create(ctx, entities.RolePermission{
		RoleId:       request.RoleId,
		PermissionId: request.PermissionId,
		Effect:       request.Effect,
	})
	end(err)
	return item, err
}

func (s *RolePermissionServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateRolePermissionRequest) (*entities.RolePermission, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.RoleId != nil {
		data["role_id"] = *request.RoleId
	}
	if request.PermissionId != nil {
		data["permission_id"] = *request.PermissionId
	}
	if request.Effect != nil {
		data["effect"] = *request.Effect
	}

	item, err := s.RolePermissionRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *RolePermissionServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.RolePermissionRepository.Delete(ctx, id)
	end(err)
	return err
}
