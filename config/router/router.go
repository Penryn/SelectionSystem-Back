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
			user.POST("/message", userController.SendConversation)
			user.GET("/message", userController.GetConversation)
			user.POST("/reason", userController.CreateReason)
			user.PUT("/reason", userController.UpdateReason)
			user.DELETE("/reason", userController.DeleteReason)
			user.GET("/reason", userController.GetReasons)
			user.POST("/post-reason", userController.PostReason)
		}
	}
}
