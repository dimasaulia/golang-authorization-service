package roles

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/roles/controllers"
)

type RoleModuleImpl struct {
	RoleController controllers.RoleController
}

func NewRoleModule(controller controllers.RoleController) *RoleModuleImpl {
	return &RoleModuleImpl{
		RoleController: controller,
	}
}

func (m *RoleModuleImpl) Name() string {
	return "roles"
}

func (m *RoleModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/roles", m.RoleController.Find)
	mux.HandleFunc("GET /api/v1/roles/by-app/{app}", m.RoleController.FindByApp)
	mux.HandleFunc("GET /api/v1/roles/{id}", m.RoleController.FindByUnique)
	mux.HandleFunc("POST /api/v1/roles", m.RoleController.Create)
	mux.HandleFunc("PUT /api/v1/roles/{id}", m.RoleController.Update)
	mux.HandleFunc("DELETE /api/v1/roles/{id}", m.RoleController.Delete)
}
