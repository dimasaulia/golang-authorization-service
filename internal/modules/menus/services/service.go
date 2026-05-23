package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/menus/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type MenuService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.Menu, error)
	FindByID(ctx context.Context, id int64) (*entities.Menu, error)
	FindByUnique(ctx context.Context, unique string) (*entities.Menu, error)
	FindByCode(ctx context.Context, code string) (*entities.Menu, error)
	FindByApp(ctx context.Context, appUnique string, params shared.ListParams) ([]entities.Menu, error)
	Create(ctx context.Context, request dto.CreateMenuRequest) (*entities.Menu, error)
	CreateBulk(ctx context.Context, request []dto.CreateMenuRequest) ([]entities.Menu, error)
	Update(ctx context.Context, id int64, request dto.UpdateMenuRequest) (*entities.Menu, error)
	Delete(ctx context.Context, id int64) error
}
