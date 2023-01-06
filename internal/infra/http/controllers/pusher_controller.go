package controllers

import (
	"chatprjkt/internal/app"
)

type PusherController struct {
	pusherService app.PusherService
}

func NewPusherController(s app.PusherService) PusherController {
	return PusherController{
		pusherService: s,
	}
}
