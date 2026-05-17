package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/users/dto"
)

type UserService interface {
	Find(ctx context.Context, limit uint64, offset uint64) ([]entities.User, error)
	FindByID(ctx context.Context, id int64) (*entities.User, error)
	Create(ctx context.Context, request dto.CreateUserRequest) (*entities.User, error)
	Update(ctx context.Context, id int64, request dto.UpdateUserRequest) (*entities.User, error)
	Delete(ctx context.Context, id int64) error
}
