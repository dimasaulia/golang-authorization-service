package controllers

import "net/http"

type AuthController interface {
	JWKS(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	KeycloakRedirect(w http.ResponseWriter, r *http.Request)
	KeycloakCallback(w http.ResponseWriter, r *http.Request)
	KeycloakExchange(w http.ResponseWriter, r *http.Request)
	GoogleRedirect(w http.ResponseWriter, r *http.Request)
	GoogleCallback(w http.ResponseWriter, r *http.Request)
	CurrentUser(w http.ResponseWriter, r *http.Request)
	CurrentUserApps(w http.ResponseWriter, r *http.Request)
	CurrentUserAccessSummary(w http.ResponseWriter, r *http.Request)
	CurrentUserAccessMenus(w http.ResponseWriter, r *http.Request)
	CurrentUserAccessPermissions(w http.ResponseWriter, r *http.Request)
	CurrentUserAccessCheck(w http.ResponseWriter, r *http.Request)
	CurrentUserAccessToken(w http.ResponseWriter, r *http.Request)
	UserApps(w http.ResponseWriter, r *http.Request)
	AccessSummary(w http.ResponseWriter, r *http.Request)
	AccessMenus(w http.ResponseWriter, r *http.Request)
	AccessPermissions(w http.ResponseWriter, r *http.Request)
	AccessCheck(w http.ResponseWriter, r *http.Request)
	AccessToken(w http.ResponseWriter, r *http.Request)
}
