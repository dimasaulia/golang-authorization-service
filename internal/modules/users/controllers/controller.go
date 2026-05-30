package controllers

import "net/http"

type UserController interface {
	Find(w http.ResponseWriter, r *http.Request)
	FindByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Signup(w http.ResponseWriter, r *http.Request)
	SignupWithGoogle(w http.ResponseWriter, r *http.Request)
	VerifyEmail(w http.ResponseWriter, r *http.Request)
	SetupPassword(w http.ResponseWriter, r *http.Request)
	ResendVerificationEmail(w http.ResponseWriter, r *http.Request)
	ResendInvitation(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
