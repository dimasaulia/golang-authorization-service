package organizations

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/organizations/controllers"
)

type OrganizationModuleImpl struct {
	OrganizationController controllers.OrganizationController
}

func NewOrganizationModule(controller controllers.OrganizationController) *OrganizationModuleImpl {
	return &OrganizationModuleImpl{
		OrganizationController: controller,
	}
}

func (m *OrganizationModuleImpl) Name() string {
	return "organizations"
}

func (m *OrganizationModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/organizations", m.OrganizationController.Find)
	mux.HandleFunc("GET /api/v1/organizations/{id}", m.OrganizationController.FindByID)
	mux.HandleFunc("POST /api/v1/organizations", m.OrganizationController.Create)
	mux.HandleFunc("PUT /api/v1/organizations/{id}", m.OrganizationController.Update)
	mux.HandleFunc("DELETE /api/v1/organizations/{id}", m.OrganizationController.Delete)
}
