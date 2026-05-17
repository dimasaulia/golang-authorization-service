package menus

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/menus/controllers"
)

type MenuModuleImpl struct {
	MenuController controllers.MenuController
}

func NewMenuModule(controller controllers.MenuController) *MenuModuleImpl {
	return &MenuModuleImpl{
		MenuController: controller,
	}
}

func (m *MenuModuleImpl) Name() string {
	return "menus"
}

func (m *MenuModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/menus", m.MenuController.Find)
	mux.HandleFunc("GET /api/v1/menus/{id}", m.MenuController.FindByID)
	mux.HandleFunc("POST /api/v1/menus", m.MenuController.Create)
	mux.HandleFunc("PUT /api/v1/menus/{id}", m.MenuController.Update)
	mux.HandleFunc("DELETE /api/v1/menus/{id}", m.MenuController.Delete)
}
