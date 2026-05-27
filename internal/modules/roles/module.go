package roles

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/roles/controllers"
	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/shared/middleware"
	sharedpermissions "github.com/open-suite/authorization/internal/shared/permissions"
)

type RoleModuleImpl struct {
	RoleController controllers.RoleController
	auth           *middleware.Authenticator
	appCode        string
}

func NewRoleModule(controller controllers.RoleController, cfg config.Config, auth *middleware.Authenticator) *RoleModuleImpl {
	return &RoleModuleImpl{
		RoleController: controller,
		auth:           auth,
		appCode:        cfg.Authz.AppCode,
	}
}

func (m *RoleModuleImpl) Name() string {
	return "roles"
}

func (m *RoleModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/roles", m.protect(sharedpermissions.AuthorizationCenterRolesRead, m.RoleController.Find))
	mux.Handle("GET /api/v1/roles/by-app/{app}", m.protect(sharedpermissions.AuthorizationCenterRolesRead, m.RoleController.FindByApp))
	mux.Handle("GET /api/v1/roles/{id}", m.protect(sharedpermissions.AuthorizationCenterRolesRead, m.RoleController.FindByUnique))
	mux.Handle("POST /api/v1/roles", m.protect(sharedpermissions.AuthorizationCenterRolesWrite, m.RoleController.Create))
	mux.Handle("PUT /api/v1/roles/{id}", m.protect(sharedpermissions.AuthorizationCenterRolesUpdate, m.RoleController.Update))
	mux.Handle("DELETE /api/v1/roles/{id}", m.protect(sharedpermissions.AuthorizationCenterRolesDelete, m.RoleController.Delete))
}

func (m *RoleModuleImpl) protect(permission string, handler http.HandlerFunc) http.Handler {
	return m.auth.RequirePermission(m.appCode, permission)(handler)
}
