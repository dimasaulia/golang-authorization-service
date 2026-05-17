package permissions

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/permissions/controllers"
)

type PermissionModuleImpl struct {
	PermissionController controllers.PermissionController
}

func NewPermissionModule(controller controllers.PermissionController) *PermissionModuleImpl {
	return &PermissionModuleImpl{
		PermissionController: controller,
	}
}

func (m *PermissionModuleImpl) Name() string {
	return "permissions"
}

func (m *PermissionModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/permissions", m.PermissionController.Find)
	mux.HandleFunc("GET /api/v1/permissions/{id}", m.PermissionController.FindByID)
	mux.HandleFunc("POST /api/v1/permissions", m.PermissionController.Create)
	mux.HandleFunc("PUT /api/v1/permissions/{id}", m.PermissionController.Update)
	mux.HandleFunc("DELETE /api/v1/permissions/{id}", m.PermissionController.Delete)
}
