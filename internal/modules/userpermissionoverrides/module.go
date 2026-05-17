package userpermissionoverrides

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/userpermissionoverrides/controllers"
)

type UserPermissionOverrideModuleImpl struct {
	UserPermissionOverrideController controllers.UserPermissionOverrideController
}

func NewUserPermissionOverrideModule(controller controllers.UserPermissionOverrideController) *UserPermissionOverrideModuleImpl {
	return &UserPermissionOverrideModuleImpl{
		UserPermissionOverrideController: controller,
	}
}

func (m *UserPermissionOverrideModuleImpl) Name() string {
	return "user_permission_overrides"
}

func (m *UserPermissionOverrideModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/user-permission-overrides", m.UserPermissionOverrideController.Find)
	mux.HandleFunc("GET /api/v1/user-permission-overrides/{id}", m.UserPermissionOverrideController.FindByID)
	mux.HandleFunc("POST /api/v1/user-permission-overrides", m.UserPermissionOverrideController.Create)
	mux.HandleFunc("PUT /api/v1/user-permission-overrides/{id}", m.UserPermissionOverrideController.Update)
	mux.HandleFunc("DELETE /api/v1/user-permission-overrides/{id}", m.UserPermissionOverrideController.Delete)
}
