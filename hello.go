package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	db := InitDB()
	defer CloseDB() // Ensure DB connection is closed on exit
	db.Close()
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		//password := c.PostForm("password")
		c.SetCookie("user", username, 3600, "/", "localhost", false, true)
		c.Redirect(http.StatusFound, "/")

	})

	r.GET("/", func(c *gin.Context) {
		username, err := c.Cookie("user")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
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
