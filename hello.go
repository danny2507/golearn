// hello.go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	InitDB()
	// Close db when finish
	defer CloseDB()

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
	r.POST("/login", login)
	r.POST("/register", register)

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

// listProducts - List all products
func listProducts(c *gin.Context) {
	products, err := ListProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// getProduct - Get a product by ID
func getProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}
	product, err := GetProduct(id)
	if err != nil || product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// addProduct - Add a new product
func addProduct(c *gin.Context) {
	var product ProductData
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	productID, err := AddProduct(product.Name, product.Description, product.Price, product.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": productID})
}

// updateProduct - Update an existing product
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
	if err := UpdateProduct(id, product.Name, product.Description, product.Price); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "product updated"})
}

// deleteProduct - Delete a product by ID
func deleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}
	if err := DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "product deleted"})
}

// login - Login by email & password
func login(c *gin.Context) {
	var requestData LoginRequestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing email/password"})
		return
	}
	email := requestData.Email
	password := requestData.Password
	userID, err := LoginUser(email, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token
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

	// Return the token
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func register(c *gin.Context) {
	var requestData RegisterRequestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	err := RegisterUser(requestData.Username, requestData.Email, requestData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "User registered"})
}
