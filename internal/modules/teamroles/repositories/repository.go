package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type TeamRoleRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.TeamRole, error)
	FindByID(ctx context.Context, id int64) (*entities.TeamRole, error)
	Create(ctx context.Context, entity entities.TeamRole) (*entities.TeamRole, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.TeamRole, error)
	Delete(ctx context.Context, id int64) error
}
