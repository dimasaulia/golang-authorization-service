package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/users/dto"
	"github.com/open-suite/authorization/internal/modules/users/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
)

type UserServiceImpl struct {
	UserRepository repositories.UserRepository
	log            *logger.LayerLogger
}

func NewUserService(repository repositories.UserRepository, appLogger *logger.Logger) UserService {
	return &UserServiceImpl{
		UserRepository: repository,
		log:            appLogger.Layer("service.users"),
	}
}

func (s *UserServiceImpl) Find(ctx context.Context, limit uint64, offset uint64) ([]entities.User, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.UserRepository.Find(ctx, limit, offset)
	end(err, "count", len(items))
	return items, err
}

func (s *UserServiceImpl) FindByID(ctx context.Context, id int64) (*entities.User, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.UserRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *UserServiceImpl) Create(ctx context.Context, request dto.CreateUserRequest) (*entities.User, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.UserRepository.Create(ctx, entities.User{
		OrganizationId: request.OrganizationId,
		Username:       request.Username,
		Email:          request.Email,
		DisplayName:    request.DisplayName,
		Type:           request.Type,
		Status:         request.Status,
	})
	end(err)
	return item, err
}

func (s *UserServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateUserRequest) (*entities.User, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.OrganizationId != nil {
		data["organization_id"] = *request.OrganizationId
	}
	if request.Username != nil {
		data["username"] = *request.Username
	}
	if request.Email != nil {
		data["email"] = *request.Email
	}
	if request.DisplayName != nil {
		data["display_name"] = *request.DisplayName
	}
	if request.Type != nil {
		data["type"] = *request.Type
	}
	if request.Status != nil {
		data["status"] = *request.Status
	}

	item, err := s.UserRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *UserServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.UserRepository.Delete(ctx, id)
	end(err)
	return err
}
