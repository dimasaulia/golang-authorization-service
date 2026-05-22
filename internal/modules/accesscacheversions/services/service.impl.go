package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/accesscacheversions/dto"
	"github.com/open-suite/authorization/internal/modules/accesscacheversions/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
)

type AccessCacheVersionServiceImpl struct {
	AccessCacheVersionRepository repositories.AccessCacheVersionRepository
	log                          *logger.LayerLogger
}

func NewAccessCacheVersionService(repository repositories.AccessCacheVersionRepository, appLogger *logger.Logger) AccessCacheVersionService {
	return &AccessCacheVersionServiceImpl{
		AccessCacheVersionRepository: repository,
		log:                          appLogger.Layer("service.access_cache_versions"),
	}
}

func (s *AccessCacheVersionServiceImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.AccessCacheVersion, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.AccessCacheVersionRepository.Find(ctx, params)
	end(err, "count", len(items))
	return items, err
}

func (s *AccessCacheVersionServiceImpl) FindByID(ctx context.Context, id int64) (*entities.AccessCacheVersion, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.AccessCacheVersionRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *AccessCacheVersionServiceImpl) Create(ctx context.Context, request dto.CreateAccessCacheVersionRequest) (*entities.AccessCacheVersion, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.AccessCacheVersionRepository.Create(ctx, entities.AccessCacheVersion{
		UserId:  request.UserId,
		AppId:   request.AppId,
		Version: request.Version,
	})
	end(err)
	return item, err
}

func (s *AccessCacheVersionServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateAccessCacheVersionRequest) (*entities.AccessCacheVersion, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.UserId != nil {
		data["user_id"] = *request.UserId
	}
	if request.AppId != nil {
		data["app_id"] = *request.AppId
	}
	if request.Version != nil {
		data["version"] = *request.Version
	}

	item, err := s.AccessCacheVersionRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *AccessCacheVersionServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.AccessCacheVersionRepository.Delete(ctx, id)
	end(err)
	return err
}
