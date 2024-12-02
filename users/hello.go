// hello.go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golearn/users/config"
	"golearn/users/controllers"
	"golearn/users/models"
	"golearn/users/services"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

var db *gorm.DB
var sdb *gorm.DB
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db = config.InitDB()
	sdb = config.InitSharedDB()
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()
	// Initialize services and controllers
	userService := services.NewUserService(db)
	sharedPostgreService := &services.SharedPostgreService{Db: sdb}
	userController := controllers.NewUserController(userService, sharedPostgreService)
	// Auto-migrate models
	db.AutoMigrate(&models.User{}, &models.Product{})

	r := gin.Default()

	// protect routes
	authorized := r.Group("/")
	authorized.Use(AuthMiddleware())

	// user API endpoints
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/register", userController.Register)
		userRoutes.POST("/login", userController.Login)
	}

	r.GET("/", func(c *gin.Context) {
		username, err := c.Cookie("user")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{"username": username})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8081") // listen and serve on 0.0.0.0:8080
}
