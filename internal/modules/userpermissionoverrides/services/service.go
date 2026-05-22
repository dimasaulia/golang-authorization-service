package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/userpermissionoverrides/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type UserPermissionOverrideService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.UserPermissionOverride, error)
	FindByID(ctx context.Context, id int64) (*entities.UserPermissionOverride, error)
	Create(ctx context.Context, request dto.CreateUserPermissionOverrideRequest) (*entities.UserPermissionOverride, error)
	Update(ctx context.Context, id int64, request dto.UpdateUserPermissionOverrideRequest) (*entities.UserPermissionOverride, error)
	Delete(ctx context.Context, id int64) error
}
