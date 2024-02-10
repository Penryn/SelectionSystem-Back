package userController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/userService"
	"SelectionSystem-Back/app/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

type CreateReasonData struct {
	UserID int    `json:"user_id"`
	Reason string `json:"reason"`
}

func CreateReason(c *gin.Context) {
	var data CreateReasonData
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
	if user.Type != 2&&user.Type != 3{
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//创建原因
	err = userService.CreateReason(user.ID, data.Reason)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, "创建成功")
}

type UpdateReasonData struct {
	ReasonID int    `json:"reason_id"`
	Reason   string `json:"reason"`
}

func UpdateReason(c *gin.Context) {
	var data UpdateReasonData
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
	if user.Type != 2&&user.Type != 3{
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//查询原因是否存在
	_, err = userService.GetReasonByID(data.ReasonID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//更新原因
	err = userService.UpdateReason(user.ID, data.ReasonID, data.Reason)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, "修改成功")
}

type DeleteReasonData struct {
	ReasonID int `json:"reason_id"`
}

func DeleteReason(c *gin.Context) {
	var data DeleteReasonData
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
	if user.Type != 2&&user.Type != 3{
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//删除原因
	err = userService.DeleteReason(user.ID, data.ReasonID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, "删除成功")
}

type GetReasonResponse struct {
	ID     int    `json:"id"`
	Reason string `json:"reason"`
}

func GetReasons(c *gin.Context) {
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
	if user.Type != 2&&user.Type != 3{
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//查询原因
	reasons, err := userService.GetReasons(user.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	var response []GetReasonResponse
	for _, reason := range reasons {
		response = append(response, GetReasonResponse{ID: reason.ID, Reason: reason.ReasonName})
	}
	utils.JsonSuccessResponse(c, response)
}

type PostReasonData struct {
	StudentID string `json:"student_id"`
	ReasonID    int `json:"reason_id"`
}

func PostReason(c *gin.Context) {
	var data PostReasonData
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
	if user.Type != 2&&user.Type != 3{
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//查询学生
	student, err := userService.GetStudentByStudentID(data.StudentID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//查询理由
	reason, err := userService.GetReasonByID(data.ReasonID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//发送理由
	err = userService.PostReason(user.ID, student.UserID, reason.ReasonName)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, "发送成功")
}
