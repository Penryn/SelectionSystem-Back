package adminController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/adminService"
	"SelectionSystem-Back/app/services/userService"
	"SelectionSystem-Back/app/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type DDLData struct {
	TimeByAdmin string `json:"time_by_admin"`
	Type 	  int    `json:"type"`
}

func SetDDL(c *gin.Context) {
	var data DDLData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	firstTime, err := time.Parse(time.RFC3339, data.TimeByAdmin)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	//获取用户id
	userID, er := c.Get("ID")
	if !er {
		utils.JsonErrorResponse(c, apiException.NoThatWrong)
		return
	}
	ID, _ := userID.(int)
	//查询用户
	user, err := userService.GetUserByID(ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//鉴权
	if user.Type != 2 && user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//获取ddl
	ddl, err := userService.GetAdminDDL()
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	if data.Type == 1 {
		if firstTime.After(ddl.SecondDDL) {
			utils.JsonErrorResponse(c, apiException.TimeSetError)
			return
		}
	} else if data.Type == 2 {
		if firstTime.Before(ddl.FirstDDL) {
			utils.JsonErrorResponse(c, apiException.TimeSetError)
			return
		}
	}
	//设置ddl
	time:=firstTime.Add((-8)*time.Hour)
	err=adminService.SetDDL(time,data.Type,user.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, "设置成功")
}


