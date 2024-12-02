// hello.go
package http_api

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golearn/users/config"
	"golearn/users/http_api/controllers"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var db *gorm.DB

func Start() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	conf := config.LoadConfig()
	ctrl := controllers.NewController(conf)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	// protect routes
	authorized := r.Group("/")
	authorized.Use(ctrl.AuthMiddleware())
	{
		authorized.GET("/products", ctrl.ListProducts)
		authorized.GET("/products/:id", ctrl.GetProduct)
		authorized.POST("/products", ctrl.AddProduct)
		authorized.PUT("/products/:id", ctrl.UpdateProduct)
		authorized.DELETE("/products/:id", ctrl.DeleteProduct)

	}
	r.Run() // listen and serve on 0.0.0.0:8080
}

//
