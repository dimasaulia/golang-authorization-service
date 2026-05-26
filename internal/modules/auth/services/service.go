package services

import (
	"context"

	"github.com/open-suite/authorization/internal/modules/auth/dto"
)

type AuthService interface {
	GoogleRedirectURL(ctx context.Context, organizationID int64) (string, error)
	HandleGoogleCallback(ctx context.Context, code string, state string) (*dto.GoogleCallbackResponse, string, error)
}
