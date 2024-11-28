package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (ctrl *Controller) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}
	product, err := ctrl.PostgreService.GetProduct(id)
	if err != nil || product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}
