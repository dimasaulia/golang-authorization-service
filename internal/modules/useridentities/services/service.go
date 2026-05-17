package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/useridentities/dto"
)

type UserIdentityService interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.UserIdentity, error)
	FindByID(ctx context.Context, id int64) (*entities.UserIdentity, error)
	Create(ctx context.Context, request dto.CreateUserIdentityRequest) (*entities.UserIdentity, error)
	Update(ctx context.Context, id int64, request dto.UpdateUserIdentityRequest) (*entities.UserIdentity, error)
	Delete(ctx context.Context, id int64) error
}
