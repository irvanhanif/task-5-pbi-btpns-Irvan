package router

import (
	galleryrouter "golang/go-jwt-mux/router/gallery.router"
	userrouter "golang/go-jwt-mux/router/user.router"

	"net/http"

	"github.com/gorilla/mux"
)


func Route(r *mux.Router) *mux.Router {
	
	filesource := http.FileServer(http.Dir("./files"))
	r.PathPrefix("/asset/").Handler(http.StripPrefix("/asset/", filesource))
	
	api := r.PathPrefix("/api").Subrouter()
	
	galleryrouter.GalleryRouter(api.PathPrefix("/photos").Subrouter())
	userrouter.UserRouter(api.PathPrefix("/users").Subrouter())

	return r
}