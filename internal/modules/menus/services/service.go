package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/menus/dto"
)

type MenuService interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.Menu, error)
	FindByID(ctx context.Context, id int64) (*entities.Menu, error)
	Create(ctx context.Context, request dto.CreateMenuRequest) (*entities.Menu, error)
	Update(ctx context.Context, id int64, request dto.UpdateMenuRequest) (*entities.Menu, error)
	Delete(ctx context.Context, id int64) error
}
