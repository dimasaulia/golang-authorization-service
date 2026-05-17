package rolepermissions

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/rolepermissions/controllers"
)

type RolePermissionModuleImpl struct {
	RolePermissionController controllers.RolePermissionController
}

func NewRolePermissionModule(controller controllers.RolePermissionController) *RolePermissionModuleImpl {
	return &RolePermissionModuleImpl{
		RolePermissionController: controller,
	}
}

func (m *RolePermissionModuleImpl) Name() string {
	return "role_permissions"
}

func (m *RolePermissionModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/role-permissions", m.RolePermissionController.Find)
	mux.HandleFunc("GET /api/v1/role-permissions/{id}", m.RolePermissionController.FindByID)
	mux.HandleFunc("POST /api/v1/role-permissions", m.RolePermissionController.Create)
	mux.HandleFunc("PUT /api/v1/role-permissions/{id}", m.RolePermissionController.Update)
	mux.HandleFunc("DELETE /api/v1/role-permissions/{id}", m.RolePermissionController.Delete)
}
