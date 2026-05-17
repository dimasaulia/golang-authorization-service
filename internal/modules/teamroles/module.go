package teamroles

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/teamroles/controllers"
)

type TeamRoleModuleImpl struct {
	TeamRoleController controllers.TeamRoleController
}

func NewTeamRoleModule(controller controllers.TeamRoleController) *TeamRoleModuleImpl {
	return &TeamRoleModuleImpl{
		TeamRoleController: controller,
	}
}

func (m *TeamRoleModuleImpl) Name() string {
	return "team_roles"
}

func (m *TeamRoleModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/team-roles", m.TeamRoleController.Find)
	mux.HandleFunc("GET /api/v1/team-roles/{id}", m.TeamRoleController.FindByID)
	mux.HandleFunc("POST /api/v1/team-roles", m.TeamRoleController.Create)
	mux.HandleFunc("PUT /api/v1/team-roles/{id}", m.TeamRoleController.Update)
	mux.HandleFunc("DELETE /api/v1/team-roles/{id}", m.TeamRoleController.Delete)
}
