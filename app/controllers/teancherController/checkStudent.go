package teancherController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/teacherService"
	"SelectionSystem-Back/app/utils"
	"github.com/gin-gonic/gin"
)

// 教师审批
type CheckData struct {
	CheckFirstPost int `json:"check_firstpost" binding:"required"`
	ID             int `json:"id" binding:"required"`
}

type CheckStudent struct {
	Checks []CheckData `json:"checks" binding:"required"`
}

func CheckByTeacher(c *gin.Context) {
	var data CheckStudent
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

	user, err := teacherService.GetUserByID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	if user.Type != 2 && user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	var count = 0
	for range data.Checks {
		count++
		if count > 6 {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	}

	for _, check := range data.Checks {
		studentInfo, err := teacherService.GetStudentInfoByID(check.ID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}

		if check.CheckFirstPost == 1 {
			studentInfo.TargetStatus = 1
		} else if check.CheckFirstPost == 2 {
			studentInfo.TargetStatus = 2
		} else {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		err = teacherService.UpdateStudentInfo(check.ID, studentInfo)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	}

	utils.JsonSuccessResponse(c, nil)
}
