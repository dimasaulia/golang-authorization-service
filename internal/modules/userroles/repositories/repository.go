package repositories

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type UserRoleRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.UserRole, error)
	FindByID(ctx context.Context, id int64) (*entities.UserRole, error)
	FindByUserAndRole(ctx context.Context, userID int64, roleID int64) (*entities.UserRole, error)
	Create(ctx context.Context, entity entities.UserRole) (*entities.UserRole, error)
	Update(ctx context.Context, id int64, data map[string]any) (*entities.UserRole, error)
	Delete(ctx context.Context, id int64) error
}
