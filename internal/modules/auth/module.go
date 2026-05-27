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
	mux.HandleFunc("GET /api/v1/auth/.well-known/jwks.json", m.AuthController.JWKS)
	mux.HandleFunc("POST /api/v1/auth/login", m.AuthController.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", m.AuthController.Refresh)
	mux.HandleFunc("POST /api/v1/auth/logout", m.AuthController.Logout)
	mux.HandleFunc("GET /api/v1/auth/google/redirect", m.AuthController.GoogleRedirect)
	mux.HandleFunc("GET /api/v1/auth/google/callback", m.AuthController.GoogleCallback)
	mux.HandleFunc("GET /api/v1/auth/users/{user_id}/apps", m.AuthController.UserApps)
	mux.HandleFunc("GET /api/v1/auth/users/{user_id}/apps/{app}/access", m.AuthController.AccessSummary)
	mux.HandleFunc("GET /api/v1/auth/users/{user_id}/apps/{app}/menus", m.AuthController.AccessMenus)
	mux.HandleFunc("GET /api/v1/auth/users/{user_id}/apps/{app}/permissions", m.AuthController.AccessPermissions)
	mux.HandleFunc("GET /api/v1/auth/users/{user_id}/apps/{app}/check", m.AuthController.AccessCheck)
	mux.HandleFunc("GET /api/v1/auth/users/{user_id}/apps/{app}/token", m.AuthController.AccessToken)
}
