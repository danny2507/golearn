package controllers

import (
	"golearn/product/config"
	"golearn/product/services"
)

type Controller struct {
	PostgreService *services.PostgreService
}

func NewController(config *config.Config) *Controller {
	return &Controller{
		PostgreService: services.NewPostgreService(config),
	}
}
