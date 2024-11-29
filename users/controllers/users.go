package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golearn/users/services"
	"net/http"
	"os"
	"time"
)

type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type UserController struct {
	UserService          *services.UserService
	SharedPostgreService *services.SharedPostgreService
}

func NewUserController(userService *services.UserService, sharedPostgreService *services.SharedPostgreService) *UserController {
	return &UserController{
		UserService:          userService,
		SharedPostgreService: sharedPostgreService,
	}
}
func (u *UserController) Register(c *gin.Context) {
	var requestData RegisterRequestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := u.UserService.RegisterUser(requestData.Username, requestData.Email, requestData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "User registered"})
}

func (u *UserController) Login(c *gin.Context) {
	var requestData LoginRequestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID, err := u.UserService.LoginUser(requestData.Email, requestData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &JWTClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	u.SharedPostgreService.AddActiveToken(tokenString)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
