package userController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/app/services/teacherService"
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
	//查询信息是否为空
	if data.Message == "" {
		utils.JsonErrorResponse(c, apiException.MessageError)
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
	UserID int `form:"user_id"`
}
type GetConversationResponse struct {
	UserAName string    `json:"user_a_name"`
	UserBName string    `json:"user_b_name"`
	Message   string    `json:"message"`
	Time      time.Time `json:"time"`
}

func GetConversation(c *gin.Context) {
	var data GetConversationData
	err := c.ShouldBindQuery(&data)
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
	//获取名称
	var usera_name string
	var userb_name string
	if userA.Type == 1 {
		user, _ := userService.GetStudentByUserID(userA.ID)
		usera_name = user.Name
	} else if userA.Type == 2 {
		user, _ := userService.GetTeacherByUserID(userA.ID)
		usera_name = user.TeacherName
	} else {
		usera_name = "管理员"
	}
	if userB.Type == 1 {
		user, _ := userService.GetStudentByUserID(userB.ID)
		userb_name = user.Name
	} else if userB.Type == 2 {
		user, _ := userService.GetTeacherByUserID(userB.ID)
		userb_name = user.TeacherName
	} else {
		userb_name = "管理员"
	}

	var response []GetConversationResponse
	for _, conversation := range conversations {
		if conversation.UserAID == userA.ID {
			response = append(response, GetConversationResponse{
				UserAName: usera_name,
				UserBName: "",
				Message:   conversation.Content,
				Time:      conversation.Time,
			})
		} else {
			response = append(response, GetConversationResponse{
				UserAName: "",
				UserBName: userb_name,
				Message:   conversation.Content,
				Time:      conversation.Time,
			})
		}
	}
	utils.JsonSuccessResponse(c, response)
}

type MessagedStudent struct {
	UserID int    `json:"user_id" binding:"required"`
	Name   string `json:"name" binding:"Required"`
}

// 获取私聊过自己的用户列表
func GetMessagedStudentList(c *gin.Context) {
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

	conversations, err := teacherService.GetMessagedStudentListByUserID(user.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	var responseStudentList = make([]MessagedStudent, 0)
	var cuser models.User
	for _, conversation := range conversations {
		if conversation.UserBID == user.ID {
			cuser, err = userService.GetUserByID(conversation.UserAID)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
		}else {
			cuser, err = userService.GetUserByID(conversation.UserBID)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
		}

		var name string
		if cuser.Type == 1 {
			studentInfo, err := userService.GetStudentByUserID(cuser.ID)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
			name = studentInfo.Name
		} else if cuser.Type == 2 {
			teacherInfo, err := userService.GetTeacherByUserID(cuser.ID)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
			name = teacherInfo.TeacherName
		} else {
			name = "管理员"
		}
		response := MessagedStudent{
			UserID: cuser.ID,
			Name:   name,
		}
		responseStudentList = append(responseStudentList, response)
	}

	utils.JsonSuccessResponse(c, responseStudentList)
}
