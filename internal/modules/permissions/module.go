package permissions

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/permissions/controllers"
	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/shared/middleware"
	sharedpermissions "github.com/open-suite/authorization/internal/shared/permissions"
)

type PermissionModuleImpl struct {
	PermissionController controllers.PermissionController
	auth                 *middleware.Authenticator
	appCode              string
}

func NewPermissionModule(controller controllers.PermissionController, cfg config.Config, auth *middleware.Authenticator) *PermissionModuleImpl {
	return &PermissionModuleImpl{
		PermissionController: controller,
		auth:                 auth,
		appCode:              cfg.Authz.AppCode,
	}
}

func (m *PermissionModuleImpl) Name() string {
	return "permissions"
}

func (m *PermissionModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/permissions", m.protect(sharedpermissions.AuthorizationCenterPermissionRead, m.PermissionController.Find))
	mux.Handle("GET /api/v1/permissions/by-app/{app}", m.protect(sharedpermissions.AuthorizationCenterPermissionRead, m.PermissionController.FindByApp))
	mux.Handle("GET /api/v1/permissions/{id}", m.protect(sharedpermissions.AuthorizationCenterPermissionRead, m.PermissionController.FindByUnique))
	mux.Handle("POST /api/v1/permissions", m.protect(sharedpermissions.AuthorizationCenterPermissionWrite, m.PermissionController.Create))
	mux.Handle("POST /api/v1/permissions/bulk", m.protect(sharedpermissions.AuthorizationCenterPermissionWrite, m.PermissionController.CreateBulk))
	mux.Handle("PUT /api/v1/permissions/{id}", m.protect(sharedpermissions.AuthorizationCenterPermissionUpdate, m.PermissionController.Update))
	mux.Handle("DELETE /api/v1/permissions/{id}", m.protect(sharedpermissions.AuthorizationCenterPermissionDelete, m.PermissionController.Delete))
}

func (m *PermissionModuleImpl) protect(permission string, handler http.HandlerFunc) http.Handler {
	return m.auth.RequirePermission(m.appCode, permission)(handler)
}
