package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/users/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type UserService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.User, error)
	FindByID(ctx context.Context, id int64) (*entities.User, error)
	Create(ctx context.Context, request dto.CreateUserRequest) (*entities.User, error)
	Update(ctx context.Context, id int64, request dto.UpdateUserRequest) (*entities.User, error)
	Delete(ctx context.Context, id int64) error
}
