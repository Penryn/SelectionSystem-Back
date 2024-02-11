package adminController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/adminService"
	"SelectionSystem-Back/app/services/userService"
	"SelectionSystem-Back/app/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type GetAdviceResponse struct {
	Name        string    `json:"name"`
	Advice      string    `json:"advice"`
	CreatedTime time.Time `json:"created_time"`
}

func GetAdvice(c *gin.Context) {
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
	//获取建议
	advices, err := adminService.GetAdvices()
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	var adviceResponse []GetAdviceResponse
	for _, advice := range advices {
		student ,err:= userService.GetStudentByID(advice.UserID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		adviceResponse = append(adviceResponse, GetAdviceResponse{
			Name:        student.Name,
			Advice:      advice.Content,
			CreatedTime: advice.CreateTime,
		})
	}
	utils.JsonSuccessResponse(c, adviceResponse)
}
