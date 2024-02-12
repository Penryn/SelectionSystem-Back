package router

import (
	"SelectionSystem-Back/app/controllers/adminController"
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
			user.POST("/avatar", userController.UploadAvatar)
		}
		admin := api.Group("/admin").Use(midwares.JWTAuthMiddleware())
		{
			admin.POST("/time", adminController.SetDDL)
			admin.GET("/advice", adminController.GetAdvice)
			admin.PUT("/reset", adminController.Reset)
		}
		student := api.Group("/student").Use(midwares.JWTAuthMiddleware())
		{
			student.POST("/info", studentController.CreatePersonalInfo)
			student.GET("/info", studentController.GetStudentInfo)
			student.PUT("/info", studentController.UpdateStudentInfo)
			student.GET("/teacher", studentController.GetTeacherList)
			student.POST("choose-teacher", studentController.PostTeacher)
			student.POST("suggest", studentController.AdvicePost)
			student.GET("get-suggest", studentController.AdviceGet)
		}
		teacher := api.Group("/teacher").Use(midwares.JWTAuthMiddleware())
		{
			teacher.GET("student", teancherController.GetStudentList)
			teacher.GET("student-check", teancherController.GetCheckStudentList)
			teacher.POST("/student/post", teancherController.CheckByTeacher)
		}
	}
}
