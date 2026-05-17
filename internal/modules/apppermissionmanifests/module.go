package apppermissionmanifests

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/apppermissionmanifests/controllers"
)

type AppPermissionManifestModuleImpl struct {
	AppPermissionManifestController controllers.AppPermissionManifestController
}

func NewAppPermissionManifestModule(controller controllers.AppPermissionManifestController) *AppPermissionManifestModuleImpl {
	return &AppPermissionManifestModuleImpl{
		AppPermissionManifestController: controller,
	}
}

func (m *AppPermissionManifestModuleImpl) Name() string {
	return "app_permission_manifests"
}

func (m *AppPermissionManifestModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/app-permission-manifests", m.AppPermissionManifestController.Find)
	mux.HandleFunc("GET /api/v1/app-permission-manifests/{id}", m.AppPermissionManifestController.FindByID)
	mux.HandleFunc("POST /api/v1/app-permission-manifests", m.AppPermissionManifestController.Create)
	mux.HandleFunc("PUT /api/v1/app-permission-manifests/{id}", m.AppPermissionManifestController.Update)
	mux.HandleFunc("DELETE /api/v1/app-permission-manifests/{id}", m.AppPermissionManifestController.Delete)
}
