package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type OrganizationRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.Organization, error)
	FindByID(ctx context.Context, id int64) (*entities.Organization, error)
	Create(ctx context.Context, entity entities.Organization) (*entities.Organization, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.Organization, error)
	Delete(ctx context.Context, id int64) error
}
