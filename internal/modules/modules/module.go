package modules

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/modules/controllers"
)

type ModuleModuleImpl struct {
	ModuleController controllers.ModuleController
}

func NewModuleModule(controller controllers.ModuleController) *ModuleModuleImpl {
	return &ModuleModuleImpl{
		ModuleController: controller,
	}
}

func (m *ModuleModuleImpl) Name() string {
	return "modules"
}

func (m *ModuleModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/modules", m.ModuleController.Find)
	mux.HandleFunc("GET /api/v1/modules/by-app/{app}", m.ModuleController.FindByApp)
	mux.HandleFunc("GET /api/v1/modules/{id}", m.ModuleController.FindByUnique)
	mux.HandleFunc("POST /api/v1/modules", m.ModuleController.Create)
	mux.HandleFunc("PUT /api/v1/modules/{id}", m.ModuleController.Update)
	mux.HandleFunc("DELETE /api/v1/modules/{id}", m.ModuleController.Delete)
}
