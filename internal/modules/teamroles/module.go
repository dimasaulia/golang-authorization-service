package teamroles

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/teamroles/controllers"
	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/shared/middleware"
	sharedpermissions "github.com/open-suite/authorization/internal/shared/permissions"
)

type TeamRoleModuleImpl struct {
	TeamRoleController controllers.TeamRoleController
	auth               *middleware.Authenticator
	appCode            string
}

func NewTeamRoleModule(controller controllers.TeamRoleController, cfg config.Config, auth *middleware.Authenticator) *TeamRoleModuleImpl {
	return &TeamRoleModuleImpl{
		TeamRoleController: controller,
		auth:               auth,
		appCode:            cfg.Authz.AppCode,
	}
}

func (m *TeamRoleModuleImpl) Name() string {
	return "team_roles"
}

func (m *TeamRoleModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/team-roles", m.protect(sharedpermissions.AuthorizationCenterTeamAndRolesRead, m.TeamRoleController.Find))
	mux.Handle("GET /api/v1/team-roles/{id}", m.protect(sharedpermissions.AuthorizationCenterTeamAndRolesRead, m.TeamRoleController.FindByID))
	mux.Handle("POST /api/v1/team-roles", m.protect(sharedpermissions.AuthorizationCenterTeamAndRolesWrite, m.TeamRoleController.Create))
	mux.Handle("PUT /api/v1/team-roles/{id}", m.protect(sharedpermissions.AuthorizationCenterTeamAndRolesUpdate, m.TeamRoleController.Update))
	mux.Handle("DELETE /api/v1/team-roles/{id}", m.protect(sharedpermissions.AuthorizationCenterTeamAndRolesDelete, m.TeamRoleController.Delete))
}

func (m *TeamRoleModuleImpl) protect(permission string, handler http.HandlerFunc) http.Handler {
	return m.auth.RequirePermission(m.appCode, permission)(handler)
}
