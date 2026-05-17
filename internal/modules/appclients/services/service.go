package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/appclients/dto"
)

type AppClientService interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.AppClient, error)
	FindByID(ctx context.Context, id int64) (*entities.AppClient, error)
	Create(ctx context.Context, request dto.CreateAppClientRequest) (*entities.AppClient, error)
	Update(ctx context.Context, id int64, request dto.UpdateAppClientRequest) (*entities.AppClient, error)
	Delete(ctx context.Context, id int64) error
}
