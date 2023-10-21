package middlewares

import (
	"golang/go-jwt-mux/config"
	"golang/go-jwt-mux/helper"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				helper.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
				return
			}
		}

		// get token value
		tokenString := c.Value

		claims := &config.JWTClaim{}

		// parsing token jwt
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				helper.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
				return
			case jwt.ValidationErrorExpired:
				helper.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Token expired"})
				return
			default:
				helper.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
				return
			}
		}

		if !token.Valid {
			helper.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
			return
		}

		next.ServeHTTP(w, r)
	})
}