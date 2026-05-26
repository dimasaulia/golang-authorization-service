//go:build wireinject
// +build wireinject

package app

import (
	"context"

	"github.com/google/wire"
	"github.com/open-suite/authorization/internal/modules/accesscacheversions"
	accesscacheversionsControllers "github.com/open-suite/authorization/internal/modules/accesscacheversions/controllers"
	accesscacheversionsRepositories "github.com/open-suite/authorization/internal/modules/accesscacheversions/repositories"
	accesscacheversionsServices "github.com/open-suite/authorization/internal/modules/accesscacheversions/services"
	"github.com/open-suite/authorization/internal/modules/actions"
	actionsControllers "github.com/open-suite/authorization/internal/modules/actions/controllers"
	actionsRepositories "github.com/open-suite/authorization/internal/modules/actions/repositories"
	actionsServices "github.com/open-suite/authorization/internal/modules/actions/services"
	"github.com/open-suite/authorization/internal/modules/appclients"
	appclientsControllers "github.com/open-suite/authorization/internal/modules/appclients/controllers"
	appclientsRepositories "github.com/open-suite/authorization/internal/modules/appclients/repositories"
	appclientsServices "github.com/open-suite/authorization/internal/modules/appclients/services"
	"github.com/open-suite/authorization/internal/modules/apppermissionmanifests"
	apppermissionmanifestsControllers "github.com/open-suite/authorization/internal/modules/apppermissionmanifests/controllers"
	apppermissionmanifestsRepositories "github.com/open-suite/authorization/internal/modules/apppermissionmanifests/repositories"
	apppermissionmanifestsServices "github.com/open-suite/authorization/internal/modules/apppermissionmanifests/services"
	"github.com/open-suite/authorization/internal/modules/apps"
	appsControllers "github.com/open-suite/authorization/internal/modules/apps/controllers"
	appsRepositories "github.com/open-suite/authorization/internal/modules/apps/repositories"
	appsServices "github.com/open-suite/authorization/internal/modules/apps/services"
	"github.com/open-suite/authorization/internal/modules/auditlogs"
	auditlogsControllers "github.com/open-suite/authorization/internal/modules/auditlogs/controllers"
	auditlogsRepositories "github.com/open-suite/authorization/internal/modules/auditlogs/repositories"
	auditlogsServices "github.com/open-suite/authorization/internal/modules/auditlogs/services"
	"github.com/open-suite/authorization/internal/modules/auth"
	authControllers "github.com/open-suite/authorization/internal/modules/auth/controllers"
	authServices "github.com/open-suite/authorization/internal/modules/auth/services"
	"github.com/open-suite/authorization/internal/modules/health"
	"github.com/open-suite/authorization/internal/modules/health/controllers"
	"github.com/open-suite/authorization/internal/modules/health/repositories"
	"github.com/open-suite/authorization/internal/modules/health/services"
	"github.com/open-suite/authorization/internal/modules/menus"
	menusControllers "github.com/open-suite/authorization/internal/modules/menus/controllers"
	menusRepositories "github.com/open-suite/authorization/internal/modules/menus/repositories"
	menusServices "github.com/open-suite/authorization/internal/modules/menus/services"
	"github.com/open-suite/authorization/internal/modules/modules"
	modulesControllers "github.com/open-suite/authorization/internal/modules/modules/controllers"
	modulesRepositories "github.com/open-suite/authorization/internal/modules/modules/repositories"
	modulesServices "github.com/open-suite/authorization/internal/modules/modules/services"
	"github.com/open-suite/authorization/internal/modules/organizations"
	organizationsControllers "github.com/open-suite/authorization/internal/modules/organizations/controllers"
	organizationsRepositories "github.com/open-suite/authorization/internal/modules/organizations/repositories"
	organizationsServices "github.com/open-suite/authorization/internal/modules/organizations/services"
	"github.com/open-suite/authorization/internal/modules/permissions"
	permissionsControllers "github.com/open-suite/authorization/internal/modules/permissions/controllers"
	permissionsRepositories "github.com/open-suite/authorization/internal/modules/permissions/repositories"
	permissionsServices "github.com/open-suite/authorization/internal/modules/permissions/services"
	"github.com/open-suite/authorization/internal/modules/releasenotes"
	releaseNoteControllers "github.com/open-suite/authorization/internal/modules/releasenotes/controllers"
	releaseNoteRepositories "github.com/open-suite/authorization/internal/modules/releasenotes/repositories"
	releaseNoteServices "github.com/open-suite/authorization/internal/modules/releasenotes/services"
	"github.com/open-suite/authorization/internal/modules/rolepermissions"
	rolepermissionsControllers "github.com/open-suite/authorization/internal/modules/rolepermissions/controllers"
	rolepermissionsRepositories "github.com/open-suite/authorization/internal/modules/rolepermissions/repositories"
	rolepermissionsServices "github.com/open-suite/authorization/internal/modules/rolepermissions/services"
	"github.com/open-suite/authorization/internal/modules/roles"
	rolesControllers "github.com/open-suite/authorization/internal/modules/roles/controllers"
	rolesRepositories "github.com/open-suite/authorization/internal/modules/roles/repositories"
	rolesServices "github.com/open-suite/authorization/internal/modules/roles/services"
	"github.com/open-suite/authorization/internal/modules/teammembers"
	teammembersControllers "github.com/open-suite/authorization/internal/modules/teammembers/controllers"
	teammembersRepositories "github.com/open-suite/authorization/internal/modules/teammembers/repositories"
	teammembersServices "github.com/open-suite/authorization/internal/modules/teammembers/services"
	"github.com/open-suite/authorization/internal/modules/teamroles"
	teamrolesControllers "github.com/open-suite/authorization/internal/modules/teamroles/controllers"
	teamrolesRepositories "github.com/open-suite/authorization/internal/modules/teamroles/repositories"
	teamrolesServices "github.com/open-suite/authorization/internal/modules/teamroles/services"
	"github.com/open-suite/authorization/internal/modules/teams"
	teamsControllers "github.com/open-suite/authorization/internal/modules/teams/controllers"
	teamsRepositories "github.com/open-suite/authorization/internal/modules/teams/repositories"
	teamsServices "github.com/open-suite/authorization/internal/modules/teams/services"
	"github.com/open-suite/authorization/internal/modules/useridentities"
	useridentitiesControllers "github.com/open-suite/authorization/internal/modules/useridentities/controllers"
	useridentitiesRepositories "github.com/open-suite/authorization/internal/modules/useridentities/repositories"
	useridentitiesServices "github.com/open-suite/authorization/internal/modules/useridentities/services"
	"github.com/open-suite/authorization/internal/modules/userpermissionoverrides"
	userpermissionoverridesControllers "github.com/open-suite/authorization/internal/modules/userpermissionoverrides/controllers"
	userpermissionoverridesRepositories "github.com/open-suite/authorization/internal/modules/userpermissionoverrides/repositories"
	userpermissionoverridesServices "github.com/open-suite/authorization/internal/modules/userpermissionoverrides/services"
	"github.com/open-suite/authorization/internal/modules/userroles"
	userrolesControllers "github.com/open-suite/authorization/internal/modules/userroles/controllers"
	userrolesRepositories "github.com/open-suite/authorization/internal/modules/userroles/repositories"
	userrolesServices "github.com/open-suite/authorization/internal/modules/userroles/services"
	"github.com/open-suite/authorization/internal/modules/users"
	usersControllers "github.com/open-suite/authorization/internal/modules/users/controllers"
	usersRepositories "github.com/open-suite/authorization/internal/modules/users/repositories"
	usersServices "github.com/open-suite/authorization/internal/modules/users/services"
	"github.com/open-suite/authorization/internal/platform/config"
	"github.com/open-suite/authorization/internal/platform/database"
	"github.com/open-suite/authorization/internal/platform/i18n"
	"github.com/open-suite/authorization/internal/platform/mailer"
	"github.com/open-suite/authorization/internal/platform/redis"
	"github.com/open-suite/authorization/internal/shared/response"
)

func Initialize(ctx context.Context) (*App, error) {
	wire.Build(
		config.Load,
		ProvideLogger,
		database.New,
		redis.New,
		i18n.NewTranslator,
		mailer.New,
		response.NewSender,
		repositories.NewHealthRepository,
		services.NewHealthService,
		controllers.NewHealthController,
		health.NewHealthModule,
		authServices.NewAuthService,
		authControllers.NewAuthController,
		auth.NewAuthModule,
		organizationsRepositories.NewOrganizationRepository,
		organizationsServices.NewOrganizationService,
		organizationsControllers.NewOrganizationController,
		organizations.NewOrganizationModule,
		usersRepositories.NewUserRepository,
		usersServices.NewUserService,
		usersControllers.NewUserController,
		users.NewUserModule,
		useridentitiesRepositories.NewUserIdentityRepository,
		useridentitiesServices.NewUserIdentityService,
		useridentitiesControllers.NewUserIdentityController,
		useridentities.NewUserIdentityModule,
		appsRepositories.NewAppRepository,
		appsServices.NewAppService,
		appsControllers.NewAppController,
		apps.NewAppModule,
		appclientsRepositories.NewAppClientRepository,
		appclientsServices.NewAppClientService,
		appclientsControllers.NewAppClientController,
		appclients.NewAppClientModule,
		modulesRepositories.NewModuleRepository,
		modulesServices.NewModuleService,
		modulesControllers.NewModuleController,
		modules.NewModuleModule,
		actionsRepositories.NewActionRepository,
		actionsServices.NewActionService,
		actionsControllers.NewActionController,
		actions.NewActionModule,
		permissionsRepositories.NewPermissionRepository,
		permissionsServices.NewPermissionService,
		permissionsControllers.NewPermissionController,
		permissions.NewPermissionModule,
		menusRepositories.NewMenuRepository,
		menusServices.NewMenuService,
		menusControllers.NewMenuController,
		menus.NewMenuModule,
		rolesRepositories.NewRoleRepository,
		rolesServices.NewRoleService,
		rolesControllers.NewRoleController,
		roles.NewRoleModule,
		rolepermissionsRepositories.NewRolePermissionRepository,
		rolepermissionsServices.NewRolePermissionService,
		rolepermissionsControllers.NewRolePermissionController,
		rolepermissions.NewRolePermissionModule,
		userrolesRepositories.NewUserRoleRepository,
		userrolesServices.NewUserRoleService,
		userrolesControllers.NewUserRoleController,
		userroles.NewUserRoleModule,
		teamsRepositories.NewTeamRepository,
		teamsServices.NewTeamService,
		teamsControllers.NewTeamController,
		teams.NewTeamModule,
		teammembersRepositories.NewTeamMemberRepository,
		teammembersServices.NewTeamMemberService,
		teammembersControllers.NewTeamMemberController,
		teammembers.NewTeamMemberModule,
		teamrolesRepositories.NewTeamRoleRepository,
		teamrolesServices.NewTeamRoleService,
		teamrolesControllers.NewTeamRoleController,
		teamroles.NewTeamRoleModule,
		userpermissionoverridesRepositories.NewUserPermissionOverrideRepository,
		userpermissionoverridesServices.NewUserPermissionOverrideService,
		userpermissionoverridesControllers.NewUserPermissionOverrideController,
		userpermissionoverrides.NewUserPermissionOverrideModule,
		accesscacheversionsRepositories.NewAccessCacheVersionRepository,
		accesscacheversionsServices.NewAccessCacheVersionService,
		accesscacheversionsControllers.NewAccessCacheVersionController,
		accesscacheversions.NewAccessCacheVersionModule,
		auditlogsRepositories.NewAuditLogRepository,
		auditlogsServices.NewAuditLogService,
		auditlogsControllers.NewAuditLogController,
		auditlogs.NewAuditLogModule,
		apppermissionmanifestsRepositories.NewAppPermissionManifestRepository,
		apppermissionmanifestsServices.NewAppPermissionManifestService,
		apppermissionmanifestsControllers.NewAppPermissionManifestController,
		apppermissionmanifests.NewAppPermissionManifestModule,
		releaseNoteRepositories.NewReleaseNoteRepository,
		releaseNoteServices.NewReleaseNoteService,
		releaseNoteControllers.NewReleaseNoteController,
		releasenotes.NewReleaseNoteModule,
		ProvideModules,
		New,
	)

	return nil, nil
}
