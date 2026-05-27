package modules

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/modules/controllers"
	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/shared/middleware"
	sharedpermissions "github.com/open-suite/authorization/internal/shared/permissions"
)

type ModuleModuleImpl struct {
	ModuleController controllers.ModuleController
	auth             *middleware.Authenticator
	appCode          string
}

func NewModuleModule(controller controllers.ModuleController, cfg config.Config, auth *middleware.Authenticator) *ModuleModuleImpl {
	return &ModuleModuleImpl{
		ModuleController: controller,
		auth:             auth,
		appCode:          cfg.Authz.AppCode,
	}
}

func (m *ModuleModuleImpl) Name() string {
	return "modules"
}

func (m *ModuleModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/modules", m.protect(sharedpermissions.AuthorizationCenterModulesRead, m.ModuleController.Find))
	mux.Handle("GET /api/v1/modules/by-app/{app}", m.protect(sharedpermissions.AuthorizationCenterModulesRead, m.ModuleController.FindByApp))
	mux.Handle("GET /api/v1/modules/{id}", m.protect(sharedpermissions.AuthorizationCenterModulesRead, m.ModuleController.FindByUnique))
	mux.Handle("POST /api/v1/modules", m.protect(sharedpermissions.AuthorizationCenterModulesWrite, m.ModuleController.Create))
	mux.Handle("PUT /api/v1/modules/{id}", m.protect(sharedpermissions.AuthorizationCenterModulesUpdate, m.ModuleController.Update))
	mux.Handle("DELETE /api/v1/modules/{id}", m.protect(sharedpermissions.AuthorizationCenterModulesDelete, m.ModuleController.Delete))
}

func (m *ModuleModuleImpl) protect(permission string, handler http.HandlerFunc) http.Handler {
	return m.auth.RequirePermission(m.appCode, permission)(handler)
}
