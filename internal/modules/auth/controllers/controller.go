package controllers

import "net/http"

type AuthController interface {
	GoogleRedirect(w http.ResponseWriter, r *http.Request)
	GoogleCallback(w http.ResponseWriter, r *http.Request)
}
