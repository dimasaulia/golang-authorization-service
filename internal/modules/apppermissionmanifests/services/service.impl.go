package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/apppermissionmanifests/dto"
	"github.com/open-suite/authorization/internal/modules/apppermissionmanifests/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared"
)

type AppPermissionManifestServiceImpl struct {
	AppPermissionManifestRepository repositories.AppPermissionManifestRepository
	log                             *logger.LayerLogger
}

func NewAppPermissionManifestService(repository repositories.AppPermissionManifestRepository, appLogger *logger.Logger) AppPermissionManifestService {
	return &AppPermissionManifestServiceImpl{
		AppPermissionManifestRepository: repository,
		log:                             appLogger.Layer("service.app_permission_manifests"),
	}
}

func (s *AppPermissionManifestServiceImpl) Find(ctx context.Context, params shared.ListParams) ([]entities.AppPermissionManifest, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.AppPermissionManifestRepository.Find(ctx, params)
	end(err, "count", len(items))
	return items, err
}

func (s *AppPermissionManifestServiceImpl) FindByID(ctx context.Context, id int64) (*entities.AppPermissionManifest, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.AppPermissionManifestRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *AppPermissionManifestServiceImpl) Create(ctx context.Context, request dto.CreateAppPermissionManifestRequest) (*entities.AppPermissionManifest, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.AppPermissionManifestRepository.Create(ctx, entities.AppPermissionManifest{
		AppId:        request.AppId,
		Version:      request.Version,
		Checksum:     request.Checksum,
		ManifestJson: request.ManifestJson,
		ActivatedAt:  request.ActivatedAt,
	})
	end(err)
	return item, err
}

func (s *AppPermissionManifestServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateAppPermissionManifestRequest) (*entities.AppPermissionManifest, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.AppId != nil {
		data["app_id"] = *request.AppId
	}
	if request.Version != nil {
		data["version"] = *request.Version
	}
	if request.Checksum != nil {
		data["checksum"] = *request.Checksum
	}
	if request.ManifestJson != nil {
		data["manifest_json"] = *request.ManifestJson
	}
	if request.ActivatedAt != nil {
		data["activated_at"] = *request.ActivatedAt
	}

	item, err := s.AppPermissionManifestRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *AppPermissionManifestServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.AppPermissionManifestRepository.Delete(ctx, id)
	end(err)
	return err
}
