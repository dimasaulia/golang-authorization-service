package dto

import "github.com/open-suite/authorization/internal/entities"

type GoogleCallbackResponse struct {
	User    *entities.User `json:"user"`
	Google  GoogleUserInfo `json:"google"`
	Created bool           `json:"created"`
}

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture,omitempty"`
}
