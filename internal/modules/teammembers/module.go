package teammembers

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/teammembers/controllers"
)

type TeamMemberModuleImpl struct {
	TeamMemberController controllers.TeamMemberController
}

func NewTeamMemberModule(controller controllers.TeamMemberController) *TeamMemberModuleImpl {
	return &TeamMemberModuleImpl{
		TeamMemberController: controller,
	}
}

func (m *TeamMemberModuleImpl) Name() string {
	return "team_members"
}

func (m *TeamMemberModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/team-members", m.TeamMemberController.Find)
	mux.HandleFunc("GET /api/v1/team-members/{id}", m.TeamMemberController.FindByID)
	mux.HandleFunc("POST /api/v1/team-members", m.TeamMemberController.Create)
	mux.HandleFunc("PUT /api/v1/team-members/{id}", m.TeamMemberController.Update)
	mux.HandleFunc("DELETE /api/v1/team-members/{id}", m.TeamMemberController.Delete)
}
