package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/modules/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type ModuleService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.Module, error)
	FindByID(ctx context.Context, id int64) (*entities.Module, error)
	FindByUnique(ctx context.Context, unique string) (*entities.Module, error)
	FindByCode(ctx context.Context, code string) (*entities.Module, error)
	FindByApp(ctx context.Context, appUnique string, params shared.ListParams) ([]entities.Module, error)
	Create(ctx context.Context, request dto.CreateModuleRequest) (*entities.Module, error)
	Update(ctx context.Context, id int64, request dto.UpdateModuleRequest) (*entities.Module, error)
	Delete(ctx context.Context, id int64) error
}
