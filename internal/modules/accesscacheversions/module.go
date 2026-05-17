package accesscacheversions

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/accesscacheversions/controllers"
)

type AccessCacheVersionModuleImpl struct {
	AccessCacheVersionController controllers.AccessCacheVersionController
}

func NewAccessCacheVersionModule(controller controllers.AccessCacheVersionController) *AccessCacheVersionModuleImpl {
	return &AccessCacheVersionModuleImpl{
		AccessCacheVersionController: controller,
	}
}

func (m *AccessCacheVersionModuleImpl) Name() string {
	return "access_cache_versions"
}

func (m *AccessCacheVersionModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/access-cache-versions", m.AccessCacheVersionController.Find)
	mux.HandleFunc("GET /api/v1/access-cache-versions/{id}", m.AccessCacheVersionController.FindByID)
	mux.HandleFunc("POST /api/v1/access-cache-versions", m.AccessCacheVersionController.Create)
	mux.HandleFunc("PUT /api/v1/access-cache-versions/{id}", m.AccessCacheVersionController.Update)
	mux.HandleFunc("DELETE /api/v1/access-cache-versions/{id}", m.AccessCacheVersionController.Delete)
}
