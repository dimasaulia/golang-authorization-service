package auth

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/auth/controllers"
	"github.com/open-suite/authorization/internal/shared/middleware"
)

type AuthModuleImpl struct {
	AuthController controllers.AuthController
	auth           *middleware.Authenticator
}

func NewAuthModule(controller controllers.AuthController, auth *middleware.Authenticator) *AuthModuleImpl {
	return &AuthModuleImpl{
		AuthController: controller,
		auth:           auth,
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
	mux.Handle("GET /api/v1/auth/me", m.protect(m.AuthController.CurrentUser))
	mux.Handle("GET /api/v1/auth/me/apps", m.protect(m.AuthController.CurrentUserApps))
	mux.Handle("GET /api/v1/auth/me/apps/{app}/access", m.protect(m.AuthController.CurrentUserAccessSummary))
	mux.Handle("GET /api/v1/auth/me/apps/{app}/menus", m.protect(m.AuthController.CurrentUserAccessMenus))
	mux.Handle("GET /api/v1/auth/me/apps/{app}/permissions", m.protect(m.AuthController.CurrentUserAccessPermissions))
	mux.Handle("GET /api/v1/auth/me/apps/{app}/check", m.protect(m.AuthController.CurrentUserAccessCheck))
	mux.Handle("GET /api/v1/auth/me/apps/{app}/token", m.protect(m.AuthController.CurrentUserAccessToken))
	mux.HandleFunc("GET /api/v1/auth/users/{user_id}/apps", m.AuthController.UserApps)
	mux.HandleFunc("GET /api/v1/auth/users/{user_id}/apps/{app}/access", m.AuthController.AccessSummary)
	mux.HandleFunc("GET /api/v1/auth/users/{user_id}/apps/{app}/menus", m.AuthController.AccessMenus)
	mux.HandleFunc("GET /api/v1/auth/users/{user_id}/apps/{app}/permissions", m.AuthController.AccessPermissions)
	mux.HandleFunc("GET /api/v1/auth/users/{user_id}/apps/{app}/check", m.AuthController.AccessCheck)
	mux.HandleFunc("GET /api/v1/auth/users/{user_id}/apps/{app}/token", m.AuthController.AccessToken)
}

func (m *AuthModuleImpl) protect(handler http.HandlerFunc) http.Handler {
	return m.auth.RequireAuthenticated(handler)
}
