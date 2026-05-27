package app

import (
	"github.com/open-suite/authorization/internal/modules/accesscacheversions"
	"github.com/open-suite/authorization/internal/modules/actions"
	"github.com/open-suite/authorization/internal/modules/appclients"
	"github.com/open-suite/authorization/internal/modules/apppermissionmanifests"
	"github.com/open-suite/authorization/internal/modules/apps"
	"github.com/open-suite/authorization/internal/modules/auditlogs"
	"github.com/open-suite/authorization/internal/modules/auth"
	authServices "github.com/open-suite/authorization/internal/modules/auth/services"
	"github.com/open-suite/authorization/internal/modules/health"
	"github.com/open-suite/authorization/internal/modules/menus"
	"github.com/open-suite/authorization/internal/modules/modules"
	"github.com/open-suite/authorization/internal/modules/organizations"
	"github.com/open-suite/authorization/internal/modules/permissions"
	"github.com/open-suite/authorization/internal/modules/releasenotes"
	"github.com/open-suite/authorization/internal/modules/rolepermissions"
	"github.com/open-suite/authorization/internal/modules/roles"
	"github.com/open-suite/authorization/internal/modules/teammembers"
	"github.com/open-suite/authorization/internal/modules/teamroles"
	"github.com/open-suite/authorization/internal/modules/teams"
	"github.com/open-suite/authorization/internal/modules/useridentities"
	"github.com/open-suite/authorization/internal/modules/userpermissionoverrides"
	"github.com/open-suite/authorization/internal/modules/userroles"
	"github.com/open-suite/authorization/internal/modules/users"
	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/platform/logger"
	"github.com/open-suite/authorization/internal/shared/middleware"
)

var _ Module = (*health.HealthModuleImpl)(nil)
var _ Module = (*organizations.OrganizationModuleImpl)(nil)
var _ Module = (*users.UserModuleImpl)(nil)
var _ Module = (*useridentities.UserIdentityModuleImpl)(nil)
var _ Module = (*apps.AppModuleImpl)(nil)
var _ Module = (*appclients.AppClientModuleImpl)(nil)
var _ Module = (*modules.ModuleModuleImpl)(nil)
var _ Module = (*actions.ActionModuleImpl)(nil)
var _ Module = (*permissions.PermissionModuleImpl)(nil)
var _ Module = (*menus.MenuModuleImpl)(nil)
var _ Module = (*roles.RoleModuleImpl)(nil)
var _ Module = (*rolepermissions.RolePermissionModuleImpl)(nil)
var _ Module = (*userroles.UserRoleModuleImpl)(nil)
var _ Module = (*teams.TeamModuleImpl)(nil)
var _ Module = (*teammembers.TeamMemberModuleImpl)(nil)
var _ Module = (*teamroles.TeamRoleModuleImpl)(nil)
var _ Module = (*userpermissionoverrides.UserPermissionOverrideModuleImpl)(nil)
var _ Module = (*accesscacheversions.AccessCacheVersionModuleImpl)(nil)
var _ Module = (*auditlogs.AuditLogModuleImpl)(nil)
var _ Module = (*apppermissionmanifests.AppPermissionManifestModuleImpl)(nil)
var _ Module = (*releasenotes.ReleaseNoteModuleImpl)(nil)
var _ Module = (*auth.AuthModuleImpl)(nil)

func ProvideLogger(cfg config.Config) (*logger.Logger, error) {
	return logger.New(logger.Config{
		Level:  cfg.Logger.Level,
		LogDir: cfg.Logger.LogDir,
	})
}

func ProvidePermissionChecker(service authServices.AuthService) middleware.PermissionChecker {
	return service
}

func ProvideModules(healthModule *health.HealthModuleImpl, authModule *auth.AuthModuleImpl, organizationModule *organizations.OrganizationModuleImpl, userModule *users.UserModuleImpl, userIdentityModule *useridentities.UserIdentityModuleImpl, appModule *apps.AppModuleImpl, appClientModule *appclients.AppClientModuleImpl, moduleModule *modules.ModuleModuleImpl, actionModule *actions.ActionModuleImpl, permissionModule *permissions.PermissionModuleImpl, menuModule *menus.MenuModuleImpl, roleModule *roles.RoleModuleImpl, rolePermissionModule *rolepermissions.RolePermissionModuleImpl, userRoleModule *userroles.UserRoleModuleImpl, teamModule *teams.TeamModuleImpl, teamMemberModule *teammembers.TeamMemberModuleImpl, teamRoleModule *teamroles.TeamRoleModuleImpl, userPermissionOverrideModule *userpermissionoverrides.UserPermissionOverrideModuleImpl, accessCacheVersionModule *accesscacheversions.AccessCacheVersionModuleImpl, auditLogModule *auditlogs.AuditLogModuleImpl, appPermissionManifestModule *apppermissionmanifests.AppPermissionManifestModuleImpl, releaseNoteModule *releasenotes.ReleaseNoteModuleImpl) []Module {
	return []Module{healthModule, authModule, organizationModule, userModule, userIdentityModule, appModule, appClientModule, moduleModule, actionModule, permissionModule, menuModule, roleModule, rolePermissionModule, userRoleModule, teamModule, teamMemberModule, teamRoleModule, userPermissionOverrideModule, accessCacheVersionModule, auditLogModule, appPermissionManifestModule, releaseNoteModule}
}
