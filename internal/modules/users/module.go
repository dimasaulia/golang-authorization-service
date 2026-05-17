package users

import (
	"net/http"

	"github.com/open-suite/authorization/internal/modules/users/controllers"
)

type UserModuleImpl struct {
	UserController controllers.UserController
}

func NewUserModule(controller controllers.UserController) *UserModuleImpl {
	return &UserModuleImpl{
		UserController: controller,
	}
}

func (m *UserModuleImpl) Name() string {
	return "users"
}

func (m *UserModuleImpl) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/users", m.UserController.Find)
	mux.HandleFunc("GET /api/v1/users/{id}", m.UserController.FindByID)
	mux.HandleFunc("POST /api/v1/users", m.UserController.Create)
	mux.HandleFunc("PUT /api/v1/users/{id}", m.UserController.Update)
	mux.HandleFunc("DELETE /api/v1/users/{id}", m.UserController.Delete)
}
