package apps

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/apps/controllers"
)

type AppModuleImpl struct {
	AppController controllers.AppController
}

func NewAppModule(controller controllers.AppController) *AppModuleImpl {
	return &AppModuleImpl{
		AppController: controller,
	}
}

func (m *AppModuleImpl) Name() string {
	return "apps"
}

func (m *AppModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/apps", m.AppController.Find)
	mux.HandleFunc("GET /api/v1/apps/{id}", m.AppController.FindByID)
	mux.HandleFunc("POST /api/v1/apps", m.AppController.Create)
	mux.HandleFunc("PUT /api/v1/apps/{id}", m.AppController.Update)
	mux.HandleFunc("DELETE /api/v1/apps/{id}", m.AppController.Delete)
}
