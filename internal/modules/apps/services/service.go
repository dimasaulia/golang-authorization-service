package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/apps/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type AppService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.App, error)
	FindByID(ctx context.Context, id int64) (*entities.App, error)
	FindByUnique(ctx context.Context, unique string) (*entities.App, error)
	FindByCode(ctx context.Context, code string) (*entities.App, error)
	Create(ctx context.Context, request dto.CreateAppRequest) (*entities.App, error)
	Update(ctx context.Context, id int64, request dto.UpdateAppRequest) (*entities.App, error)
	Delete(ctx context.Context, id int64) error
}
