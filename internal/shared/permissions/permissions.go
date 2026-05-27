package permissions

const (
	AppCodeAuthorizationCenter = "authorization-center"

	ActionDelete = "delete"
	ActionUpdate = "update"
	ActionWrite  = "write"
	ActionRead   = "read"

	ModuleApps               = "apps"
	ModuleModules            = "modules"
	ModulePermission         = "permission"
	ModuleMenuAndRoutes      = "menu-and-routes"
	ModuleRoles              = "roles"
	ModuleRolesAndPermission = "roles-and-permission"
	ModuleDashboard          = "dashboard"
	ModuleUsers              = "users"
	ModuleActions            = "actions"
	ModuleTeams              = "teams"
	ModuleTeamAndRoles       = "team-and-roles"
)

const (
	AuthorizationCenterAppsDelete = AppCodeAuthorizationCenter + "." + ModuleApps + "." + ActionDelete
	AuthorizationCenterAppsUpdate = AppCodeAuthorizationCenter + "." + ModuleApps + "." + ActionUpdate
	AuthorizationCenterAppsWrite  = AppCodeAuthorizationCenter + "." + ModuleApps + "." + ActionWrite
	AuthorizationCenterAppsRead   = AppCodeAuthorizationCenter + "." + ModuleApps + "." + ActionRead

	AuthorizationCenterModulesDelete = AppCodeAuthorizationCenter + "." + ModuleModules + "." + ActionDelete
	AuthorizationCenterModulesUpdate = AppCodeAuthorizationCenter + "." + ModuleModules + "." + ActionUpdate
	AuthorizationCenterModulesWrite  = AppCodeAuthorizationCenter + "." + ModuleModules + "." + ActionWrite
	AuthorizationCenterModulesRead   = AppCodeAuthorizationCenter + "." + ModuleModules + "." + ActionRead

	AuthorizationCenterPermissionDelete = AppCodeAuthorizationCenter + "." + ModulePermission + "." + ActionDelete
	AuthorizationCenterPermissionUpdate = AppCodeAuthorizationCenter + "." + ModulePermission + "." + ActionUpdate
	AuthorizationCenterPermissionWrite  = AppCodeAuthorizationCenter + "." + ModulePermission + "." + ActionWrite
	AuthorizationCenterPermissionRead   = AppCodeAuthorizationCenter + "." + ModulePermission + "." + ActionRead

	AuthorizationCenterMenuAndRoutesDelete = AppCodeAuthorizationCenter + "." + ModuleMenuAndRoutes + "." + ActionDelete
	AuthorizationCenterMenuAndRoutesUpdate = AppCodeAuthorizationCenter + "." + ModuleMenuAndRoutes + "." + ActionUpdate
	AuthorizationCenterMenuAndRoutesWrite  = AppCodeAuthorizationCenter + "." + ModuleMenuAndRoutes + "." + ActionWrite
	AuthorizationCenterMenuAndRoutesRead   = AppCodeAuthorizationCenter + "." + ModuleMenuAndRoutes + "." + ActionRead

	AuthorizationCenterRolesDelete = AppCodeAuthorizationCenter + "." + ModuleRoles + "." + ActionDelete
	AuthorizationCenterRolesUpdate = AppCodeAuthorizationCenter + "." + ModuleRoles + "." + ActionUpdate
	AuthorizationCenterRolesWrite  = AppCodeAuthorizationCenter + "." + ModuleRoles + "." + ActionWrite
	AuthorizationCenterRolesRead   = AppCodeAuthorizationCenter + "." + ModuleRoles + "." + ActionRead

	AuthorizationCenterRolesAndPermissionDelete = AppCodeAuthorizationCenter + "." + ModuleRolesAndPermission + "." + ActionDelete
	AuthorizationCenterRolesAndPermissionUpdate = AppCodeAuthorizationCenter + "." + ModuleRolesAndPermission + "." + ActionUpdate
	AuthorizationCenterRolesAndPermissionWrite  = AppCodeAuthorizationCenter + "." + ModuleRolesAndPermission + "." + ActionWrite
	AuthorizationCenterRolesAndPermissionRead   = AppCodeAuthorizationCenter + "." + ModuleRolesAndPermission + "." + ActionRead

	AuthorizationCenterDashboardDelete = AppCodeAuthorizationCenter + "." + ModuleDashboard + "." + ActionDelete
	AuthorizationCenterDashboardUpdate = AppCodeAuthorizationCenter + "." + ModuleDashboard + "." + ActionUpdate
	AuthorizationCenterDashboardWrite  = AppCodeAuthorizationCenter + "." + ModuleDashboard + "." + ActionWrite
	AuthorizationCenterDashboardRead   = AppCodeAuthorizationCenter + "." + ModuleDashboard + "." + ActionRead

	AuthorizationCenterUsersDelete = AppCodeAuthorizationCenter + "." + ModuleUsers + "." + ActionDelete
	AuthorizationCenterUsersUpdate = AppCodeAuthorizationCenter + "." + ModuleUsers + "." + ActionUpdate
	AuthorizationCenterUsersWrite  = AppCodeAuthorizationCenter + "." + ModuleUsers + "." + ActionWrite
	AuthorizationCenterUsersRead   = AppCodeAuthorizationCenter + "." + ModuleUsers + "." + ActionRead

	AuthorizationCenterActionsDelete = AppCodeAuthorizationCenter + "." + ModuleActions + "." + ActionDelete
	AuthorizationCenterActionsUpdate = AppCodeAuthorizationCenter + "." + ModuleActions + "." + ActionUpdate
	AuthorizationCenterActionsWrite  = AppCodeAuthorizationCenter + "." + ModuleActions + "." + ActionWrite
	AuthorizationCenterActionsRead   = AppCodeAuthorizationCenter + "." + ModuleActions + "." + ActionRead

	AuthorizationCenterTeamsDelete = AppCodeAuthorizationCenter + "." + ModuleTeams + "." + ActionDelete
	AuthorizationCenterTeamsUpdate = AppCodeAuthorizationCenter + "." + ModuleTeams + "." + ActionUpdate
	AuthorizationCenterTeamsWrite  = AppCodeAuthorizationCenter + "." + ModuleTeams + "." + ActionWrite
	AuthorizationCenterTeamsRead   = AppCodeAuthorizationCenter + "." + ModuleTeams + "." + ActionRead

	AuthorizationCenterTeamAndRolesDelete = AppCodeAuthorizationCenter + "." + ModuleTeamAndRoles + "." + ActionDelete
	AuthorizationCenterTeamAndRolesUpdate = AppCodeAuthorizationCenter + "." + ModuleTeamAndRoles + "." + ActionUpdate
	AuthorizationCenterTeamAndRolesWrite  = AppCodeAuthorizationCenter + "." + ModuleTeamAndRoles + "." + ActionWrite
	AuthorizationCenterTeamAndRolesRead   = AppCodeAuthorizationCenter + "." + ModuleTeamAndRoles + "." + ActionRead
)
