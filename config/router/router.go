package router

import (
	"SelectionSystem-Back/app/controllers/studentController"
	"SelectionSystem-Back/app/controllers/teancherController"
	"SelectionSystem-Back/app/controllers/userController"
	"SelectionSystem-Back/app/midwares"
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	const pre = "/api"
	api := r.Group(pre)
	{
		api.POST("/login", userController.Login)

		user := api.Group("/user").Use(midwares.JWTAuthMiddleware())
		{
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
