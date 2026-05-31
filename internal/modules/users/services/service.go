package services

import (
	"context"

	"github.com/open-suite/authorization/internal/entities"
	"github.com/open-suite/authorization/internal/modules/users/dto"
	"github.com/open-suite/authorization/internal/shared"
)

type UserService interface {
	Find(ctx context.Context, params shared.ListParams) ([]entities.User, error)
	FindByID(ctx context.Context, id int64) (*dto.UserResponse, error)
	Create(ctx context.Context, request dto.CreateUserRequest) (*dto.UserResponse, error)
	Signup(ctx context.Context, request dto.SignupUserRequest) (*dto.UserResponse, error)
	SignupWithGoogle(ctx context.Context, request dto.GoogleSignupRequest) (*dto.UserResponse, error)
	VerifyEmail(ctx context.Context, request dto.VerifyEmailRequest) error
	SetupPassword(ctx context.Context, request dto.SetupPasswordRequest) error
	ResendVerificationEmail(ctx context.Context, request dto.ResendVerificationEmailRequest) error
	ResendInvitation(ctx context.Context, id int64, request dto.ResendInvitationRequest) error
	Update(ctx context.Context, id int64, request dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id int64) error
}
