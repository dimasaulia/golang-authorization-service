package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/actions/dto"
	"github.com/open-suite/authorization/internal/modules/actions/repositories"
	"github.com/open-suite/authorization/internal/platform/logger"
)

type ActionServiceImpl struct {
	ActionRepository repositories.ActionRepository
	log              *logger.LayerLogger
}

func NewActionService(repository repositories.ActionRepository, appLogger *logger.Logger) ActionService {
	return &ActionServiceImpl{
		ActionRepository: repository,
		log:              appLogger.Layer("service.actions"),
	}
}

func (s *ActionServiceImpl) Find(ctx context.Context, limit uint64, offset uint64) ([]entities.Action, error) {
	end := s.log.Start(ctx, "Find")
	items, err := s.ActionRepository.Find(ctx, limit, offset)
	end(err, "count", len(items))
	return items, err
}

func (s *ActionServiceImpl) FindByID(ctx context.Context, id int64) (*entities.Action, error) {
	end := s.log.Start(ctx, "FindByID", "id", id)
	item, err := s.ActionRepository.FindByID(ctx, id)
	end(err)
	return item, err
}

func (s *ActionServiceImpl) Create(ctx context.Context, request dto.CreateActionRequest) (*entities.Action, error) {
	end := s.log.Start(ctx, "Create")
	item, err := s.ActionRepository.Create(ctx, entities.Action{
		Code: request.Code,
		Name: request.Name,
	})
	end(err)
	return item, err
}

func (s *ActionServiceImpl) Update(ctx context.Context, id int64, request dto.UpdateActionRequest) (*entities.Action, error) {
	end := s.log.Start(ctx, "Update", "id", id)

	data := map[string]any{}
	if request.Code != nil {
		data["code"] = *request.Code
	}
	if request.Name != nil {
		data["name"] = *request.Name
	}

	item, err := s.ActionRepository.Update(ctx, id, data)
	end(err)
	return item, err
}

func (s *ActionServiceImpl) Delete(ctx context.Context, id int64) error {
	end := s.log.Start(ctx, "Delete", "id", id)
	err := s.ActionRepository.Delete(ctx, id)
	end(err)
	return err
}
