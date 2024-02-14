package timeController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/app/services/timeService"
	"SelectionSystem-Back/app/utils"

	"github.com/gin-gonic/gin"
)

func Allocate(c *gin.Context) {
	// 获取所有学生
	var students []models.Student
	students, err := timeService.QueryStudents()
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	for _, student := range students {
		if student.TargetID == 0 {
			//获取所有教师
			teachers_id, err := timeService.QueryTeachers()
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
			//分配教师
			student.TargetID = timeService.RandomTeacher(teachers_id)
			//更新学生信息
			err = timeService.UpdateStudent(student)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
			//更新教师信息
			err = timeService.UpdateTeacher(student.TargetID)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}

		}
	}
}
