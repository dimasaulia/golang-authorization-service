package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type TeamMemberRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.TeamMember, error)
	FindByID(ctx context.Context, id int64) (*entities.TeamMember, error)
	Create(ctx context.Context, entity entities.TeamMember) (*entities.TeamMember, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.TeamMember, error)
	Delete(ctx context.Context, id int64) error
}
