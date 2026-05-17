package useridentities

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/useridentities/controllers"
)

type UserIdentityModuleImpl struct {
	UserIdentityController controllers.UserIdentityController
}

func NewUserIdentityModule(controller controllers.UserIdentityController) *UserIdentityModuleImpl {
	return &UserIdentityModuleImpl{
		UserIdentityController: controller,
	}
}

func (m *UserIdentityModuleImpl) Name() string {
	return "user_identities"
}

func (m *UserIdentityModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/user-identities", m.UserIdentityController.Find)
	mux.HandleFunc("GET /api/v1/user-identities/{id}", m.UserIdentityController.FindByID)
	mux.HandleFunc("POST /api/v1/user-identities", m.UserIdentityController.Create)
	mux.HandleFunc("PUT /api/v1/user-identities/{id}", m.UserIdentityController.Update)
	mux.HandleFunc("DELETE /api/v1/user-identities/{id}", m.UserIdentityController.Delete)
}
