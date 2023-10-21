package galleryrouter

import (
	productcontroller "golang/go-jwt-mux/controllers/productController"
	"golang/go-jwt-mux/middlewares"

	"github.com/gorilla/mux"
)

func GalleryRouter(r *mux.Router) {
	r.HandleFunc("/", productcontroller.UploadPhotoProduct).Methods("POST")
	r.HandleFunc("/", productcontroller.GetPhotosProduct).Methods("GET")
	r.HandleFunc("/{id}", productcontroller.GetDetailPhoto).Methods("GET")
	r.HandleFunc("/{id}", productcontroller.UpdateProductGallery).Methods("PUT")
	r.HandleFunc("/{id}", productcontroller.DeleteProductGallery).Methods("DELETE")
	r.Use(middlewares.JWTMiddleware)
}