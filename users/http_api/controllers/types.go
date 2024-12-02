package controllers

import (
	"github.com/golang-jwt/jwt"
	"os"
)

type RegisterRequestData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequestData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
