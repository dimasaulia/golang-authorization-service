package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/appclients/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type AppClientService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.AppClient, error)
	FindByID(ctx context.Context, id int64) (*entities.AppClient, error)
	Create(ctx context.Context, request dto.CreateAppClientRequest) (*entities.AppClient, error)
	Update(ctx context.Context, id int64, request dto.UpdateAppClientRequest) (*entities.AppClient, error)
	Delete(ctx context.Context, id int64) error
}
