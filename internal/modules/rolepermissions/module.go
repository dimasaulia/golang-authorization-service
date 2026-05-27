package rolepermissions

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/rolepermissions/controllers"
	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/shared/middleware"
	sharedpermissions "github.com/open-suite/authorization/internal/shared/permissions"
)

type RolePermissionModuleImpl struct {
	RolePermissionController controllers.RolePermissionController
	auth                     *middleware.Authenticator
	appCode                  string
}

func NewRolePermissionModule(controller controllers.RolePermissionController, cfg config.Config, auth *middleware.Authenticator) *RolePermissionModuleImpl {
	return &RolePermissionModuleImpl{
		RolePermissionController: controller,
		auth:                     auth,
		appCode:                  cfg.Authz.AppCode,
	}
}

func (m *RolePermissionModuleImpl) Name() string {
	return "role_permissions"
}

func (m *RolePermissionModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/v1/role-permissions", m.protect(sharedpermissions.AuthorizationCenterRolesAndPermissionRead, m.RolePermissionController.Find))
	mux.Handle("GET /api/v1/role-permissions/by-app/{app}", m.protect(sharedpermissions.AuthorizationCenterRolesAndPermissionRead, m.RolePermissionController.FindByApp))
	mux.Handle("GET /api/v1/role-permissions/by-role/{role}", m.protect(sharedpermissions.AuthorizationCenterRolesAndPermissionRead, m.RolePermissionController.FindByRole))
	mux.Handle("GET /api/v1/role-permissions/available-permissions/by-app/{app}", m.protect(sharedpermissions.AuthorizationCenterRolesAndPermissionRead, m.RolePermissionController.FindAvailablePermissionsByApp))
	mux.Handle("GET /api/v1/role-permissions/roles", m.protect(sharedpermissions.AuthorizationCenterRolesAndPermissionRead, m.RolePermissionController.FindRoleSummaries))
	mux.Handle("GET /api/v1/role-permissions/roles/by-app/{app}", m.protect(sharedpermissions.AuthorizationCenterRolesAndPermissionRead, m.RolePermissionController.FindRoleSummariesByApp))
	mux.Handle("GET /api/v1/role-permissions/{id}", m.protect(sharedpermissions.AuthorizationCenterRolesAndPermissionRead, m.RolePermissionController.FindByID))
	mux.Handle("POST /api/v1/role-permissions", m.protect(sharedpermissions.AuthorizationCenterRolesAndPermissionWrite, m.RolePermissionController.Create))
	mux.Handle("POST /api/v1/role-permissions/bulk", m.protect(sharedpermissions.AuthorizationCenterRolesAndPermissionWrite, m.RolePermissionController.CreateBulk))
	mux.Handle("PUT /api/v1/role-permissions/{id}", m.protect(sharedpermissions.AuthorizationCenterRolesAndPermissionUpdate, m.RolePermissionController.Update))
	mux.Handle("PUT /api/v1/role-permissions/by-role/{role}", m.protect(sharedpermissions.AuthorizationCenterRolesAndPermissionUpdate, m.RolePermissionController.UpdateByRole))
	mux.Handle("DELETE /api/v1/role-permissions/{id}", m.protect(sharedpermissions.AuthorizationCenterRolesAndPermissionDelete, m.RolePermissionController.Delete))
}

func (m *RolePermissionModuleImpl) protect(permission string, handler http.HandlerFunc) http.Handler {
	return m.auth.RequirePermission(m.appCode, permission)(handler)
}
