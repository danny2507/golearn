package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (ctrl *Controller) ListProducts(c *gin.Context) {
	products, err := ctrl.PostgreService.ListProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve products"})
		return
	}
	c.JSON(http.StatusOK, products)
}
