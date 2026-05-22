package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/useridentities/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type UserIdentityService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.UserIdentity, error)
	FindByID(ctx context.Context, id int64) (*entities.UserIdentity, error)
	Create(ctx context.Context, request dto.CreateUserIdentityRequest) (*entities.UserIdentity, error)
	Update(ctx context.Context, id int64, request dto.UpdateUserIdentityRequest) (*entities.UserIdentity, error)
	Delete(ctx context.Context, id int64) error
}
