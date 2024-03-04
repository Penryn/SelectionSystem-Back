package adminController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/adminService"
	"SelectionSystem-Back/app/services/userService"
	"SelectionSystem-Back/app/utils"
	"math"

	"github.com/gin-gonic/gin"
)

type ResetData struct {
	UserID int `json:"user_id"`
}

func ResetUser(c *gin.Context) {
	var data ResetData
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
	//鉴权
	if user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//重置密码
	err = adminService.ResetPassword(data.UserID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

type GetUserByAdminData struct {
	PageNum  int    `form:"page_num" validate:"required"`
	PageSize int    `form:"page_size" validate:"required"`
	UserName string `form:"user_name"`
}

type GetUserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Type     int    `json:"type"`
}

func GetUserByAdmin(c *gin.Context) {
	var data GetUserByAdminData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	userID, er := c.Get("ID")
	if !er {
		utils.JsonErrorResponse(c, apiException.NoThatWrong)
		return
	}
	ID, _ := userID.(int)
	user, err := userService.GetUserByID(ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	if user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	var num *int64
	users, num, err := adminService.GetUsers(data.PageNum, data.PageSize, data.UserName)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	response := make([]GetUserResponse, 0)
	for _, v := range users {
		response = append(response, GetUserResponse{
			ID:       v.ID,
			Username: v.Username,
			Type:     v.Type,
		})
	}
	utils.JsonSuccessResponse(c, gin.H{
		"data":           response,
		"total_page_num": math.Ceil(float64(*num) / float64(data.PageSize)),
		"user_num":       *num,
	})
}
