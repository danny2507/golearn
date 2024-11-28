package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (ctrl *Controller) AddProduct(c *gin.Context) {
	var product ProductData
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	productID, err := ctrl.PostgreService.AddProduct(product.Name, product.Description, product.Price, product.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": productID})
}
