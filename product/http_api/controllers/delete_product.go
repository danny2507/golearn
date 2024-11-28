package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (ctrl *Controller) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}
	if err := ctrl.PostgreService.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "product deleted"})
}
