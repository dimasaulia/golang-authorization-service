package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type TeamRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.Team, error)
	FindByID(ctx context.Context, id int64) (*entities.Team, error)
	FindByCode(ctx context.Context, code string) (*entities.Team, error)
	Create(ctx context.Context, entity entities.Team) (*entities.Team, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.Team, error)
	Delete(ctx context.Context, id int64) error
}
