package studentController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/studentService"
	"SelectionSystem-Back/app/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type PostTeacherID struct {
	TargetID int `json:"teacher_id" binding:"required"`
}

// 提交目标导师
func PostTeacher(c *gin.Context) {
	var data PostTeacherID
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}

	userId, er := c.Get("UserID")
	if !er {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	adminDDL, err := studentService.GetAdminDDL()
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	currentTime := time.Now()
	if currentTime.After(adminDDL.FirstDDL) {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	targetTeacher, err := studentService.GetTeacherByTeacherID(data.TargetID)
	if err != nil && err != gorm.ErrRecordNotFound {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	teacherDDL, err := studentService.GetTeacherDDLByUserID(targetTeacher.UserID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	if currentTime.After(teacherDDL.FirstDDL) {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	if targetTeacher.StudentsNum >= 6 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	studentInfo, err := studentService.GetStudentInfoByUserID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	err = studentService.UpdateTargetTeacher(userId.(int), data.TargetID, studentInfo)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}
