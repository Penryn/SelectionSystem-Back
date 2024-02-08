package router

import (
	"SelectionSystem-Back/app/controllers/userController"

	"github.com/gin-gonic/gin"
	"SelectionSystem-Back/app/midwares"
)

func Init(r *gin.Engine) {
	const pre = "/api"
	api := r.Group(pre)
	{
		api.POST("/login", userController.Login)
		user:=api.Group("/user").Use(midwares.JWTAuthMiddleware())
		{
			user.PUT("/reset", userController.Retrieve)
		}
	}
}
