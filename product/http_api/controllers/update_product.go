package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (ctrl *Controller) UpdateProduct(c *gin.Context) {
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
	if err := ctrl.PostgreService.UpdateProduct(id, product.Name, product.Description, product.Price); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "product updated"})
}
