package controllers

import "net/http"

type RolePermissionController interface {
	Find(w http.ResponseWriter, r *http.Request)
	FindByApp(w http.ResponseWriter, r *http.Request)
	FindByRole(w http.ResponseWriter, r *http.Request)
	FindRoleSummaries(w http.ResponseWriter, r *http.Request)
	FindRoleSummariesByApp(w http.ResponseWriter, r *http.Request)
	FindAvailablePermissionsByApp(w http.ResponseWriter, r *http.Request)
	FindByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	CreateBulk(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	UpdateByRole(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
