package main

import (
	"golang/go-jwt-mux/models"
	"golang/go-jwt-mux/router"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	models.ConnectDatabase()

	r := router.Route(mux.NewRouter())
	
	log.Fatal(http.ListenAndServe(":8080", r))
}