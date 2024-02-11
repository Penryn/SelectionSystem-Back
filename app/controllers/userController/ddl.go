package userController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/userService"
	"SelectionSystem-Back/app/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type DDLResponse struct {
	FirstTime  time.Time `json:"first_time_by_admin"`
	SecondTime time.Time `json:"second_time_by_admin"`
}

func GetAdminDDL(c *gin.Context) {
	//获取ddl
	ddl, err := userService.GetAdminDDL()
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, DDLResponse{FirstTime: ddl.FirstDDL, SecondTime: ddl.SecondDDL})
}


type DDLData struct {
	TeacherID    int    `form:"teacher_id"`
}

func GetTeacherDDL(c *gin.Context) {
	var data DDLData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	//获取教师
	teacher, err := userService.GetTeacherByTeacherID(data.TeacherID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//获取ddl
	ddl, err := userService.GetTeacherDDLTime(teacher.UserID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, DDLResponse{FirstTime: ddl.FirstDDL})

}

