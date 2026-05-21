package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/actions/dto"
)

type ActionService interface {
	Find(ctx context.Context, limit uint64, offset uint64, search string) ([]entities.Action, error)
	FindByID(ctx context.Context, id int64) (*entities.Action, error)
	Create(ctx context.Context, request dto.CreateActionRequest) (*entities.Action, error)
	Update(ctx context.Context, id int64, request dto.UpdateActionRequest) (*entities.Action, error)
	Delete(ctx context.Context, id int64) error
}
