package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/userpermissionoverrides/dto"
)

type UserPermissionOverrideService interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.UserPermissionOverride, error)
	FindByID(ctx context.Context, id int64) (*entities.UserPermissionOverride, error)
	Create(ctx context.Context, request dto.CreateUserPermissionOverrideRequest) (*entities.UserPermissionOverride, error)
	Update(ctx context.Context, id int64, request dto.UpdateUserPermissionOverrideRequest) (*entities.UserPermissionOverride, error)
	Delete(ctx context.Context, id int64) error
}
