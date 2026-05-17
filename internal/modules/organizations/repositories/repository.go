package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
)

type OrganizationRepository interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.Organization, error)
	FindByID(ctx context.Context, id int64) (*entities.Organization, error)
	Create(ctx context.Context, entity entities.Organization) (*entities.Organization, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.Organization, error)
	Delete(ctx context.Context, id int64) error
}
