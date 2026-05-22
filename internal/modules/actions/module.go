package actions

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/actions/controllers"
)

type ActionModuleImpl struct {
	ActionController controllers.ActionController
}

func NewActionModule(controller controllers.ActionController) *ActionModuleImpl {
	return &ActionModuleImpl{
		ActionController: controller,
	}
}

func (m *ActionModuleImpl) Name() string {
	return "actions"
}

func (m *ActionModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/actions", m.ActionController.Find)
	mux.HandleFunc("GET /api/v1/actions/{id}", m.ActionController.FindByUnique)
	mux.HandleFunc("POST /api/v1/actions", m.ActionController.Create)
	mux.HandleFunc("PUT /api/v1/actions/{id}", m.ActionController.Update)
	mux.HandleFunc("DELETE /api/v1/actions/{id}", m.ActionController.Delete)
}
