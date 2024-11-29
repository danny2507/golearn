package controllers

import (
	"golearn/product/services"
	"gorm.io/gorm"
)

type Controller struct {
	PostgreService       *services.PostgreService
	SharedPostgreService *services.SharedPostgreService
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{
		PostgreService: services.NewPostgreService(db),
	}
}
