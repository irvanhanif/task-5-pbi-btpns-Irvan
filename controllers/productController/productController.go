package productcontroller

import (
	"fmt"
	"golang/go-jwt-mux/helper"
	"golang/go-jwt-mux/models"
	"net/http"

	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func UploadPhotoProduct(w http.ResponseWriter, r *http.Request)  {
	if err := r.ParseMultipartForm(1024); err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}
	
	// alg penyimpanan foto
	success, code, message := helper.FileUploaded(r)
	if !success {
		helper.ResponseJSON(w, code, map[string]string{"message": message})
		return
	}

	// mengambil input body
	var fileInput models.Photo

	r.ParseForm()
	fileInput.Title = r.FormValue("Title")
	fileInput.Caption = r.FormValue("Caption")

	fileInput.PhotoUrl = r.Host + "/asset/" + message
	
	c, _ := r.Cookie("token")
	token := c.Value
	// parsing token jwt
	tokenVal, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	var username string
	if tokenClaims, ok := tokenVal.Claims.(jwt.MapClaims); ok {
		username = fmt.Sprint(tokenClaims["Username"])
	}

	var user models.User
	if err := models.DB.Where("username = ?", username).First(&user).Error; err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error(), "user": username})
		return
	}

	fileInput.UserID = user.Id
	fileInput.User = user

	// insert to db
	if err := models.DB.Create(&fileInput).Error; err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	// response, _ := json.Marshal(fileInput)
	helper.ResponseJSON(w, http.StatusOK, map[string]string{"message": "Success upload"})
	// helper.ResponseJSON(w, http.StatusOK, map[string]string{"message": username})
}

func GetPhotosProduct(w http.ResponseWriter, r *http.Request)  {
	var photos []models.Photo

	if err := models.DB.Preload("User").Find(&photos).Error; err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	helper.ResponseJSON(w, http.StatusOK, photos)
}

func GetDetailPhoto(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	var photo models.Photo
	if err := models.DB.Preload("User").First(&photo, id).Error ;err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helper.ResponseJSON(w, http.StatusNotFound, map[string]string{"message": "Photo not found"})
			return
		default:
			helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}
	}
	
	helper.ResponseJSON(w, http.StatusOK, photo)
}

func modelsDbId(id int) *gorm.DB {
	return models.DB.Where("id = ?", id)
}

func UpdateProductGallery(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	var photo models.Photo
	defer r.Body.Close()

	if err := models.DB.First(&photo, id).Error ;err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helper.ResponseJSON(w, http.StatusNotFound, map[string]string{"message": "Photo not found"})
			return
		default:
			helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}
	}

	isPhotoUploaded := true
	
	uploadedFile, _, err := r.FormFile("file")
	if err != nil {
		if err.Error() == "http: no such file" {
			isPhotoUploaded = false
		} else{
			helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}
	}
	
	if isPhotoUploaded && uploadedFile != nil {
	// if isPhotoUploaded {
		// delete photo exist in db
		if err := os.Remove("./files/" + photo.PhotoUrl); err != nil {
			helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}

		// alg penyimpanan foto
		success, code, message := helper.FileUploaded(r)
		if !success {
			helper.ResponseJSON(w, code, map[string]string{"message": message})
			return
		}

		photo.PhotoUrl = message
	}

	// update data
	// mengambil input body
	r.ParseForm()
	photo.Title = r.FormValue("Title")
	photo.Caption = r.FormValue("Caption")

	if err := modelsDbId(int(id)).Updates(&photo).Error; err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	helper.ResponseJSON(w, http.StatusOK, map[string]string{"message": "Success updated"})
}

func DeleteProductGallery(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}

	var photo models.Photo
	if err := modelsDbId(int(id)).First(&photo).Error; err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	if err := os.Remove("./files/" + photo.PhotoUrl); err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	if modelsDbId(int(id)).Delete(&photo).RowsAffected == 0 {
		helper.ResponseJSON(w, http.StatusBadRequest, map[string]string{"message": "Photo cant be deleted"})
		return
	}

	helper.ResponseJSON(w, http.StatusOK, map[string]string{"message": "Success deleted"})
}