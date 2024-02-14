package userController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/userService"
	"SelectionSystem-Back/app/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type SendConversationData struct {
	UserBID int    `json:"user_b_id"`
	Message string `json:"message"`
}

func SendConversation(c *gin.Context) {
	var data SendConversationData
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
	//查询用户a
	userA, err := userService.GetUserByID(ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//查询用户b
	userB, err := userService.GetUserByID(data.UserBID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//发送消息
	err = userService.SendConversation(userA.ID, userB.ID, data.Message)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, "发送成功")
}

type GetConversationData struct {
	UserID int `json:"user_id"`
}
type GetConversationResponse struct {
	UserAName string    `json:"user_a_name"`
	UserBName string    `json:"user_b_name"`
	Message   string    `json:"message"`
	Time      time.Time `json:"time"`
}

func GetConversation(c *gin.Context) {
	var data GetConversationData
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
	//查询用户a
	userA, err := userService.GetUserByID(ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//查询用户b
	userB, err := userService.GetUserByID(data.UserID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//获取消息
	conversations, err := userService.GetConversation(userA.ID, userB.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	var response []GetConversationResponse
	for _, conversation := range conversations {
		if conversation.UserAID == userA.ID {
			response = append(response, GetConversationResponse{
				UserAName: userA.Username,
				UserBName: userB.Username,
				Message:   conversation.Content,
				Time:      conversation.Time,
			})
		} else {
			response = append(response, GetConversationResponse{
				UserAName: userB.Username,
				UserBName: userA.Username,
				Message:   conversation.Content,
				Time:      conversation.Time,
			})
		}
	}
	utils.JsonSuccessResponse(c, response)
}
