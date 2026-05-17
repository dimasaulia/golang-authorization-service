package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/userpermissionoverrides/dto"
	"github.com/open-suite/authorization/internal/modules/userpermissionoverrides/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
)

type UserPermissionOverrideServiceImpl struct {
	UserPermissionOverrideRepository repositories.UserPermissionOverrideRepository
	log                              *logger.LayerLogger
}

func NewUserPermissionOverrideService(repository repositories.UserPermissionOverrideRepository, appLogger *logger.Logger) UserPermissionOverrideService {
	return &UserPermissionOverrideServiceImpl{
		UserPermissionOverrideRepository: repository,
		log:                              appLogger.Layer("service.user_permission_overrides"),
	}
}

func (s *UserPermissionOverrideServiceImpl) Find(ctx context.Context, limit uint64, offset uint64) ([]entities.UserPermissionOverride, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.UserPermissionOverrideRepository.Find(ctx, limit, offset)
	end(err, "count", len(items))
	return items, err
}

func (s *UserPermissionOverrideServiceImpl) FindByID(ctx context.Context, id int64) (*entities.UserPermissionOverride, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.UserPermissionOverrideRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *UserPermissionOverrideServiceImpl) Create(ctx context.Context, request dto.CreateUserPermissionOverrideRequest) (*entities.UserPermissionOverride, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.UserPermissionOverrideRepository.Create(ctx, entities.UserPermissionOverride{
		UserId:       request.UserId,
		PermissionId: request.PermissionId,
		Effect:       request.Effect,
		Reason:       request.Reason,
		ExpiresAt:    request.ExpiresAt,
		CreatedBy:    request.CreatedBy,
	})
	end(err)
	return item, err
}

func (s *UserPermissionOverrideServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateUserPermissionOverrideRequest) (*entities.UserPermissionOverride, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.UserId != nil {
		data["user_id"] = *request.UserId
	}
	if request.PermissionId != nil {
		data["permission_id"] = *request.PermissionId
	}
	if request.Effect != nil {
		data["effect"] = *request.Effect
	}
	if request.Reason != nil {
		data["reason"] = *request.Reason
	}
	if request.ExpiresAt != nil {
		data["expires_at"] = *request.ExpiresAt
	}
	if request.CreatedBy != nil {
		data["created_by"] = *request.CreatedBy
	}

	item, err := s.UserPermissionOverrideRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *UserPermissionOverrideServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.UserPermissionOverrideRepository.Delete(ctx, id)
	end(err)
	return err
}
