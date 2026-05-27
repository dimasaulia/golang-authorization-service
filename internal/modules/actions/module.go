package actions

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/actions/controllers"
	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/shared/middleware"
	sharedpermissions "github.com/open-suite/authorization/internal/shared/permissions"
)

type ActionModuleImpl struct {
	ActionController controllers.ActionController
	auth             *middleware.Authenticator
	appCode          string
}

func NewActionModule(controller controllers.ActionController, cfg config.Config, auth *middleware.Authenticator) *ActionModuleImpl {
	return &ActionModuleImpl{
		ActionController: controller,
		auth:             auth,
		appCode:          cfg.Authz.AppCode,
	}
}

func (m *ActionModuleImpl) Name() string {
	return "actions"
}

func (m *ActionModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/actions", m.protect(sharedpermissions.AuthorizationCenterActionsRead, m.ActionController.Find))
	mux.Handle("GET /api/v1/actions/{id}", m.protect(sharedpermissions.AuthorizationCenterActionsRead, m.ActionController.FindByUnique))
	mux.Handle("POST /api/v1/actions", m.protect(sharedpermissions.AuthorizationCenterActionsWrite, m.ActionController.Create))
	mux.Handle("PUT /api/v1/actions/{id}", m.protect(sharedpermissions.AuthorizationCenterActionsUpdate, m.ActionController.Update))
	mux.Handle("DELETE /api/v1/actions/{id}", m.protect(sharedpermissions.AuthorizationCenterActionsDelete, m.ActionController.Delete))
}

func (m *ActionModuleImpl) protect(permission string, handler http.HandlerFunc) http.Handler {
	return m.auth.RequirePermission(m.appCode, permission)(handler)
}
