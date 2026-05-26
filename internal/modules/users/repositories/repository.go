package repositories

import (
	"context"
	"time"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/shared"
)

type CreateUserInput struct {
	User               entities.User
	PasswordHash       string
	MustChangePassword bool
	Identity           *entities.UserIdentity
	VerificationCode   *CreateVerificationCodeInput
}

type CreateVerificationCodeInput struct {
	Purpose   string
	CodeHash  string
	ExpiresAt time.Time
}

type UserRepository interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.User, error)
	FindByID(ctx context.Context, id int64) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByIdentity(ctx context.Context, provider string, providerUserID string) (*entities.User, error)
	Create(ctx context.Context, input CreateUserInput) (*entities.User, error)
	LinkIdentity(ctx context.Context, identity entities.UserIdentity) error
	Update(ctx context.Context, id int64, data map[string]any) (*entities.User, error)
	Delete(ctx context.Context, id int64) error
	FindVerificationCode(ctx context.Context, purpose string, codeHash string) (*entities.UserVerificationCode, error)
	UseVerificationCode(ctx context.Context, codeID int64) error
	MarkEmailVerified(ctx context.Context, userID int64) error
}
