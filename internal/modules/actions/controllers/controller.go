package controllers

import "net/http"

type ActionController interface {
	Find(w http.ResponseWriter, r *http.Request)
	FindByUnique(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
