package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
)

type TeamRoleRepository interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.TeamRole, error)
	FindByID(ctx context.Context, id int64) (*entities.TeamRole, error)
	Create(ctx context.Context, entity entities.TeamRole) (*entities.TeamRole, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.TeamRole, error)
	Delete(ctx context.Context, id int64) error
}
