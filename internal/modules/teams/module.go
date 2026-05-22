package teams

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/teams/controllers"
)

type TeamModuleImpl struct {
	TeamController controllers.TeamController
}

func NewTeamModule(controller controllers.TeamController) *TeamModuleImpl {
	return &TeamModuleImpl{
		TeamController: controller,
	}
}

func (m *TeamModuleImpl) Name() string {
	return "teams"
}

func (m *TeamModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/teams", m.TeamController.Find)
	mux.HandleFunc("GET /api/v1/teams/{id}", m.TeamController.FindByUnique)
	mux.HandleFunc("POST /api/v1/teams", m.TeamController.Create)
	mux.HandleFunc("PUT /api/v1/teams/{id}", m.TeamController.Update)
	mux.HandleFunc("DELETE /api/v1/teams/{id}", m.TeamController.Delete)
}
