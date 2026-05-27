package teams

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/teams/controllers"
	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/shared/middleware"
	sharedpermissions "github.com/open-suite/authorization/internal/shared/permissions"
)

type TeamModuleImpl struct {
	TeamController controllers.TeamController
	auth           *middleware.Authenticator
	appCode        string
}

func NewTeamModule(controller controllers.TeamController, cfg config.Config, auth *middleware.Authenticator) *TeamModuleImpl {
	return &TeamModuleImpl{
		TeamController: controller,
		auth:           auth,
		appCode:        cfg.Authz.AppCode,
	}
}

func (m *TeamModuleImpl) Name() string {
	return "teams"
}

func (m *TeamModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/teams", m.protect(sharedpermissions.AuthorizationCenterTeamsRead, m.TeamController.Find))
	mux.Handle("GET /api/v1/teams/{id}", m.protect(sharedpermissions.AuthorizationCenterTeamsRead, m.TeamController.FindByUnique))
	mux.Handle("POST /api/v1/teams", m.protect(sharedpermissions.AuthorizationCenterTeamsWrite, m.TeamController.Create))
	mux.Handle("PUT /api/v1/teams/{id}", m.protect(sharedpermissions.AuthorizationCenterTeamsUpdate, m.TeamController.Update))
	mux.Handle("DELETE /api/v1/teams/{id}", m.protect(sharedpermissions.AuthorizationCenterTeamsDelete, m.TeamController.Delete))
}

func (m *TeamModuleImpl) protect(permission string, handler http.HandlerFunc) http.Handler {
	return m.auth.RequirePermission(m.appCode, permission)(handler)
}
