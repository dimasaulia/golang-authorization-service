package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/organizations/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type OrganizationService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.Organization, error)
	FindByID(ctx context.Context, id int64) (*entities.Organization, error)
	Create(ctx context.Context, request dto.CreateOrganizationRequest) (*entities.Organization, error)
	Update(ctx context.Context, id int64, request dto.UpdateOrganizationRequest) (*entities.Organization, error)
	Delete(ctx context.Context, id int64) error
}
