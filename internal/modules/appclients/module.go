package appclients

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/appclients/controllers"
)

type AppClientModuleImpl struct {
	AppClientController controllers.AppClientController
}

func NewAppClientModule(controller controllers.AppClientController) *AppClientModuleImpl {
	return &AppClientModuleImpl{
		AppClientController: controller,
	}
}

func (m *AppClientModuleImpl) Name() string {
	return "app_clients"
}

func (m *AppClientModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/app-clients", m.AppClientController.Find)
	mux.HandleFunc("GET /api/v1/app-clients/{id}", m.AppClientController.FindByID)
	mux.HandleFunc("POST /api/v1/app-clients", m.AppClientController.Create)
	mux.HandleFunc("PUT /api/v1/app-clients/{id}", m.AppClientController.Update)
	mux.HandleFunc("DELETE /api/v1/app-clients/{id}", m.AppClientController.Delete)
}
