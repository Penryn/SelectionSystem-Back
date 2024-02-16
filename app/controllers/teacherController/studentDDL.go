package teacherController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/teacherService"
	"SelectionSystem-Back/app/utils"
	"github.com/gin-gonic/gin"
	"time"
)

type DDLSetData struct {
	TimeByTeacher string `json:"time_by_teacher" binding:"required"`
}

func DDLSetByTeacher(c *gin.Context) {
	var data DDLSetData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}

	userId, er := c.Get("ID")
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

	ddlTime, err := time.Parse(time.RFC3339, data.TimeByTeacher)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	adminDDL, err := teacherService.GetAdminDDL()
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	if ddlTime.After(adminDDL.FirstDDL) {
		ddlTime = adminDDL.FirstDDL
	} else {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}

	err = teacherService.SetDDL(ddlTime, userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}
