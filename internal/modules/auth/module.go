package auth

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/auth/controllers"
)

type AuthModuleImpl struct {
	AuthController controllers.AuthController
}

func NewAuthModule(controller controllers.AuthController) *AuthModuleImpl {
	return &AuthModuleImpl{
		AuthController: controller,
	}
}

func (m *AuthModuleImpl) Name() string {
	return "auth"
}

func (m *AuthModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/auth/google/redirect", m.AuthController.GoogleRedirect)
	mux.HandleFunc("GET /api/v1/auth/google/callback", m.AuthController.GoogleCallback)
}
