package controllers

import "net/http"

type AuthController interface {
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	GoogleRedirect(w http.ResponseWriter, r *http.Request)
	GoogleCallback(w http.ResponseWriter, r *http.Request)
	UserApps(w http.ResponseWriter, r *http.Request)
	AccessSummary(w http.ResponseWriter, r *http.Request)
	AccessMenus(w http.ResponseWriter, r *http.Request)
	AccessPermissions(w http.ResponseWriter, r *http.Request)
	AccessCheck(w http.ResponseWriter, r *http.Request)
	AccessToken(w http.ResponseWriter, r *http.Request)
}
