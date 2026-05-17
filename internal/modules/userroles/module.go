package userroles

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/userroles/controllers"
)

type UserRoleModuleImpl struct {
	UserRoleController controllers.UserRoleController
}

func NewUserRoleModule(controller controllers.UserRoleController) *UserRoleModuleImpl {
	return &UserRoleModuleImpl{
		UserRoleController: controller,
	}
}

func (m *UserRoleModuleImpl) Name() string {
	return "user_roles"
}

func (m *UserRoleModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/user-roles", m.UserRoleController.Find)
	mux.HandleFunc("GET /api/v1/user-roles/{id}", m.UserRoleController.FindByID)
	mux.HandleFunc("POST /api/v1/user-roles", m.UserRoleController.Create)
	mux.HandleFunc("PUT /api/v1/user-roles/{id}", m.UserRoleController.Update)
	mux.HandleFunc("DELETE /api/v1/user-roles/{id}", m.UserRoleController.Delete)
}
