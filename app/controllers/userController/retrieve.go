package userController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/userService"
	"SelectionSystem-Back/app/utils"

	"github.com/gin-gonic/gin"
)

type RetrieveData struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func Retrieve(c *gin.Context) {
	var data RetrieveData
	err := c.ShouldBindJSON(&data)
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
	//判断旧密码是否正确
	if user.Password != data.OldPassword {
		utils.JsonErrorResponse(c, apiException.NoThatWrong)
		return
	}
	//修改密码
	err = userService.UpdatePassword(user, data.NewPassword)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, "修改成功")

}
