package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("kdjsahkjdkwjqjo13029139pdk")

type JWTClaim struct {
	Username string
	jwt.RegisteredClaims
}