package controllers

import (
	"golearn/users/config"
	"golearn/users/services"
)

type Controller struct {
	PostgreService *services.PostgreService
}

func NewController(config *config.Config) *Controller {
	return &Controller{
		PostgreService: services.NewPostgreService(config),
	}
}
