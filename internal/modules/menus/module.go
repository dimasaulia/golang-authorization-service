package menus

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/menus/controllers"
	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/shared/middleware"
	sharedpermissions "github.com/open-suite/authorization/internal/shared/permissions"
)

type MenuModuleImpl struct {
	MenuController controllers.MenuController
	auth           *middleware.Authenticator
	appCode        string
}

func NewMenuModule(controller controllers.MenuController, cfg config.Config, auth *middleware.Authenticator) *MenuModuleImpl {
	return &MenuModuleImpl{
		MenuController: controller,
		auth:           auth,
		appCode:        cfg.Authz.AppCode,
	}
}

func (m *MenuModuleImpl) Name() string {
	return "menus"
}

func (m *MenuModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/menus", m.protect(sharedpermissions.AuthorizationCenterMenuAndRoutesRead, m.MenuController.Find))
	mux.Handle("GET /api/v1/menus/by-app/{app}", m.protect(sharedpermissions.AuthorizationCenterMenuAndRoutesRead, m.MenuController.FindByApp))
	mux.Handle("GET /api/v1/menus/{id}", m.protect(sharedpermissions.AuthorizationCenterMenuAndRoutesRead, m.MenuController.FindByUnique))
	mux.Handle("POST /api/v1/menus", m.protect(sharedpermissions.AuthorizationCenterMenuAndRoutesWrite, m.MenuController.Create))
	mux.Handle("POST /api/v1/menus/bulk", m.protect(sharedpermissions.AuthorizationCenterMenuAndRoutesWrite, m.MenuController.CreateBulk))
	mux.Handle("PUT /api/v1/menus/{id}", m.protect(sharedpermissions.AuthorizationCenterMenuAndRoutesUpdate, m.MenuController.Update))
	mux.Handle("DELETE /api/v1/menus/{id}", m.protect(sharedpermissions.AuthorizationCenterMenuAndRoutesDelete, m.MenuController.Delete))
}

func (m *MenuModuleImpl) protect(permission string, handler http.HandlerFunc) http.Handler {
	return m.auth.RequirePermission(m.appCode, permission)(handler)
}
