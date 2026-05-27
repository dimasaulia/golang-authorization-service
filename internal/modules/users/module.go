package users

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/users/controllers"
	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/shared/middleware"
	sharedpermissions "github.com/open-suite/authorization/internal/shared/permissions"
)

type UserModuleImpl struct {
	UserController controllers.UserController
	auth           *middleware.Authenticator
	appCode        string
}

func NewUserModule(controller controllers.UserController, cfg config.Config, auth *middleware.Authenticator) *UserModuleImpl {
	return &UserModuleImpl{
		UserController: controller,
		auth:           auth,
		appCode:        cfg.Authz.AppCode,
	}
}

func (m *UserModuleImpl) Name() string {
	return "users"
}

func (m *UserModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/users", m.protect(sharedpermissions.AuthorizationCenterUsersRead, m.UserController.Find))
	mux.HandleFunc("POST /api/v1/users/signup", m.UserController.Signup)
	mux.HandleFunc("POST /api/v1/users/signup/google", m.UserController.SignupWithGoogle)
	mux.HandleFunc("GET /api/v1/users/verify-email", m.UserController.VerifyEmail)
	mux.HandleFunc("POST /api/v1/users/verify-email", m.UserController.VerifyEmail)
	mux.Handle("GET /api/v1/users/{id}", m.protect(sharedpermissions.AuthorizationCenterUsersRead, m.UserController.FindByID))
	mux.Handle("POST /api/v1/users", m.protect(sharedpermissions.AuthorizationCenterUsersWrite, m.UserController.Create))
	mux.Handle("PUT /api/v1/users/{id}", m.protect(sharedpermissions.AuthorizationCenterUsersUpdate, m.UserController.Update))
	mux.Handle("DELETE /api/v1/users/{id}", m.protect(sharedpermissions.AuthorizationCenterUsersDelete, m.UserController.Delete))
}

func (m *UserModuleImpl) protect(permission string, handler http.HandlerFunc) http.Handler {
	return m.auth.RequirePermission(m.appCode, permission)(handler)
}
