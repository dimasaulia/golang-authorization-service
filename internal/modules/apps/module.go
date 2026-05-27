package apps

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/apps/controllers"
	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/shared/middleware"
	sharedpermissions "github.com/open-suite/authorization/internal/shared/permissions"
)

type AppModuleImpl struct {
	AppController controllers.AppController
	auth          *middleware.Authenticator
	appCode       string
}

func NewAppModule(controller controllers.AppController, cfg config.Config, auth *middleware.Authenticator) *AppModuleImpl {
	return &AppModuleImpl{
		AppController: controller,
		auth:          auth,
		appCode:       cfg.Authz.AppCode,
	}
}

func (m *AppModuleImpl) Name() string {
	return "apps"
}

func (m *AppModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/apps", m.protect(sharedpermissions.AuthorizationCenterAppsRead, m.AppController.Find))
	mux.Handle("GET /api/v1/apps/{id}", m.protect(sharedpermissions.AuthorizationCenterAppsRead, m.AppController.FindByUnique))
	mux.Handle("POST /api/v1/apps", m.protect(sharedpermissions.AuthorizationCenterAppsWrite, m.AppController.Create))
	mux.Handle("PUT /api/v1/apps/{id}", m.protect(sharedpermissions.AuthorizationCenterAppsUpdate, m.AppController.Update))
	mux.Handle("DELETE /api/v1/apps/{id}", m.protect(sharedpermissions.AuthorizationCenterAppsDelete, m.AppController.Delete))
}

func (m *AppModuleImpl) protect(permission string, handler http.HandlerFunc) http.Handler {
	return m.auth.RequirePermission(m.appCode, permission)(handler)
}
