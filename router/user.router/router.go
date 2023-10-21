package userrouter

import (
	authcontroller "golang/go-jwt-mux/controllers/authController"

	"github.com/gorilla/mux"
)

func UserRouter(r *mux.Router) {
	r.HandleFunc("/login", authcontroller.Login).Methods("POST")
	r.HandleFunc("/register", authcontroller.Register).Methods("POST")
	r.HandleFunc("/logout", authcontroller.Logout).Methods("DELETE")
}