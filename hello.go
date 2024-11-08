package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"strconv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	InitDB()
	defer CloseDB()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		c.SetCookie("user", username, 3600, "/", "localhost", false, true)
		c.Redirect(http.StatusFound, "/")
	})

	// Product API Endpoints
	r.GET("/products", listProducts)
	r.GET("/products/:id", getProduct)
	r.POST("/products", addProduct)
	r.PUT("/products/:id", updateProduct)
	r.DELETE("/products/:id", deleteProduct)

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
	var product Product
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
	var product Product
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
