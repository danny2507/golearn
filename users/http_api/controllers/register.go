package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (u *Controller) Register(c *gin.Context) {
	var requestData RegisterRequestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := u.PostgreService.RegisterUser(requestData.Username, requestData.Email, requestData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "User registered"})
}
