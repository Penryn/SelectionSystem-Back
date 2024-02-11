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
	TimeByAdmin string `json:"first_time_by_admin"`
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
	//设置ddl
	err=adminService.SetDDL(firstTime,data.Type,user.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, "设置成功")
}


