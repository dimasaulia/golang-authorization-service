package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/useridentities/dto"
	"github.com/open-suite/authorization/internal/modules/useridentities/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
)

type UserIdentityServiceImpl struct {
	UserIdentityRepository repositories.UserIdentityRepository
	log                    *logger.LayerLogger
}

func NewUserIdentityService(repository repositories.UserIdentityRepository, appLogger *logger.Logger) UserIdentityService {
	return &UserIdentityServiceImpl{
		UserIdentityRepository: repository,
		log:                    appLogger.Layer("service.user_identities"),
	}
}

func (s *UserIdentityServiceImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.UserIdentity, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.UserIdentityRepository.Find(ctx, params)
	end(err, "count", len(items))
	return items, err
}

func (s *UserIdentityServiceImpl) FindByID(ctx context.Context, id int64) (*entities.UserIdentity, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.UserIdentityRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *UserIdentityServiceImpl) Create(ctx context.Context, request dto.CreateUserIdentityRequest) (*entities.UserIdentity, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.UserIdentityRepository.Create(ctx, entities.UserIdentity{
		UserId:         request.UserId,
		Provider:       request.Provider,
		ProviderUserId: request.ProviderUserId,
		Username:       request.Username,
		Email:          request.Email,
		IsPrimary:      request.IsPrimary,
	})
	end(err)
	return item, err
}

func (s *UserIdentityServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateUserIdentityRequest) (*entities.UserIdentity, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.UserId != nil {
		data["user_id"] = *request.UserId
	}
	if request.Provider != nil {
		data["provider"] = *request.Provider
	}
	if request.ProviderUserId != nil {
		data["provider_user_id"] = *request.ProviderUserId
	}
	if request.Username != nil {
		data["username"] = *request.Username
	}
	if request.Email != nil {
		data["email"] = *request.Email
	}
	if request.IsPrimary != nil {
		data["is_primary"] = *request.IsPrimary
	}

	item, err := s.UserIdentityRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *UserIdentityServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.UserIdentityRepository.Delete(ctx, id)
	end(err)
	return err
}
