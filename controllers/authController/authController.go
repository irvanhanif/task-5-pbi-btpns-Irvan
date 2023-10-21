package authcontroller

import (
	"encoding/json"
	"golang/go-jwt-mux/config"
	"golang/go-jwt-mux/helper"
	"golang/go-jwt-mux/models"
	"net/http"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// mengambil input body
	var userInput models.UserAuth
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	defer r.Body.Close()

	pattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if userInput.Email == "" || userInput.Password == "" {
		helper.ResponseJSON(w, http.StatusNotAcceptable, map[string]string{"message": "email atau password kosong"})
		return
	}else if !pattern.MatchString(userInput.Email) {
		helper.ResponseJSON(w, http.StatusNotAcceptable, map[string]string{"message": "email is not correct"})
		return
	}

	// get data user from db 
	var user models.User
	if err := models.DB.Where("email = ?", userInput.Email).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			helper.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Username atau Password salah"})
			return
		default:
			helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
			return
		}
	}
	
	// cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		helper.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Username atau Password salah"})
		return
	}

	// set jwt
	expTime := time.Now().Add(time.Minute * 60)
	claims := &config.JWTClaim{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "go-jwt",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	tokenAlg := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenAlg.SignedString(config.JWT_KEY)
	if err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	// set token in cookie
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Path: "/",
		Value: token,
		HttpOnly: true,
	})

	helper.ResponseJSON(w, http.StatusOK, map[string]string{"message": "Login berhasil"})
	return
}

func Register(w http.ResponseWriter, r *http.Request) {
	// mengambil input body
	var userInput models.UserAuth
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		helper.ResponseJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		return
	}
	defer r.Body.Close()

	var userExist models.User
	pattern := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if userInput.Email == "" || userInput.Password == "" {
		helper.ResponseJSON(w, http.StatusNotAcceptable, map[string]string{"message": "email atau password kosong"})
		return
	}else if !pattern.MatchString(userInput.Email) {
		helper.ResponseJSON(w, http.StatusNotAcceptable, map[string]string{"message": "email is not correct"})
		return
	}else if models.DB.Where("email = ?", userInput.Email).First(&userExist); userExist.Id > 0 {
		helper.ResponseJSON(w, http.StatusNotAcceptable, map[string]string{"message": "user is exist"})
		return
	}
	
	// hash pass by bcrypt
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashPassword)
	
	var userData models.User
	userData.Username = userInput.Username
	userData.Email = userInput.Email
	userData.Password = userInput.Password
	
	// insert to db
	if err := models.DB.Create(&userData).Error; err != nil {
		helper.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	helper.ResponseJSON(w, http.StatusOK, map[string]string{"message": "success"})
	return
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// del token in cookie
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Path: "/",
		Value: "",
		HttpOnly: true,
		MaxAge: -1,
	})

	helper.ResponseJSON(w, http.StatusOK, map[string]string{"message": "Logout berhasil"})
	return
}