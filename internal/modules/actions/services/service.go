package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/actions/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type ActionService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.Action, error)
	FindByID(ctx context.Context, id int64) (*entities.Action, error)
	FindByUnique(ctx context.Context, unique string) (*entities.Action, error)
	FindByCode(ctx context.Context, code string) (*entities.Action, error)
	Create(ctx context.Context, request dto.CreateActionRequest) (*entities.Action, error)
	Update(ctx context.Context, id int64, request dto.UpdateActionRequest) (*entities.Action, error)
	Delete(ctx context.Context, id int64) error
}
