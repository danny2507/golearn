// hello.go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golearn/config"
	"golearn/controllers"
	"golearn/models"
	"golearn/services"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
)

type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

var db *gorm.DB
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db = config.InitDB()
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}()
	// Initialize services and controllers
	userService := services.NewUserService(db)
	userController := controllers.NewUserController(userService)
	// Auto-migrate models
	db.AutoMigrate(&models.User{}, &models.Product{})

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	// protect routes
	authorized := r.Group("/")
	authorized.Use(AuthMiddleware())
	{
		authorized.GET("/products", listProducts)
		authorized.GET("/products/:id", getProduct)
		authorized.POST("/products", addProduct)
		authorized.PUT("/products/:id", updateProduct)
		authorized.DELETE("/products/:id", deleteProduct)

	}
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

	r.Run() // listen and serve on 0.0.0.0:8080
}

// Handler functions

func listProducts(c *gin.Context) {
	products, err := services.ListProducts(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

func getProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}
	product, err := services.GetProduct(db, id)
	if err != nil || product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func addProduct(c *gin.Context) {
	var product ProductData
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	productID, err := services.AddProduct(db, product.Name, product.Description, product.Price, product.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": productID})
}

func updateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}
	var product ProductData
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := services.UpdateProduct(db, id, product.Name, product.Description, product.Price); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "product updated"})
}

func deleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}
	if err := services.DeleteProduct(db, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "product deleted"})
}
