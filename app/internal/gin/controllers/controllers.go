package controllers

import (
	"github.com/gin-gonic/gin"
	"itmo-profile/internal/gin/handlers"
)

func SetupControllers(r *gin.Engine) {
	api := r.Group("profile/api-v1")
	{
		// NOT AUTH
		api.POST("/login", handlers.Login)
		api.POST("/be-itmo", handlers.BeItmo)
		api.POST("/be-itmo/vote", handlers.Vote)
	}
}
