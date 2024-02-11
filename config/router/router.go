package router

import (
	"SelectionSystem-Back/app/controllers/userController"

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
			user.POST("avatar", studentController.UploadAvatar)
		}

		student := api.Group("/student").Use(midwares.JWTAuthMiddleware())
		{
			student.POST("/info", studentController.CreatePersonalInfo)
			student.GET("/info", studentController.GetStudentInfo)
			student.PUT("/info", studentController.UpdateStudentInfo)
			student.GET("/teacher", studentController.GetTeacherList)
		}

		teacher := api.Group("/teacher").Use(midwares.JWTAuthMiddleware())
		{
			teacher.GET("student", teancherController.GetStudentList)
			teacher.GET("student-check", teancherController.GetCheckStudentList)
		}
	}
}
