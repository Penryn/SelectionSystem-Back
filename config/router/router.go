package router

import (
	"SelectionSystem-Back/app/controllers/adminController"
	"SelectionSystem-Back/app/controllers/userController"

	"SelectionSystem-Back/app/midwares"

	"github.com/gin-gonic/gin"
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
			user.GET("/admin/time", userController.GetAdminDDL)
			user.GET("/teacher/time", userController.GetTeacherDDL)
		}
		admin:=api.Group("/admin").Use(midwares.JWTAuthMiddleware())
		{
			admin.POST("/time", adminController.SetDDL)
			admin.GET("/advice", adminController.GetAdvice)
			admin.PUT("/reset", adminController.Reset)
		}
	}
}
