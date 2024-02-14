package adminController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/adminService"
	"SelectionSystem-Back/app/services/userService"
	"SelectionSystem-Back/app/utils"
	"math"
	"time"

	"github.com/gin-gonic/gin"
)

type GetAdviceData struct {
	PageNum  int `form:"page_num" validate:"required"`
	PageSize int `form:"page_size" validate:"required"`
}

type GetAdviceResponse struct {
	Name        string    `json:"name"`
	Advice      string    `json:"advice"`
	CreatedTime time.Time `json:"created_time"`
}

func GetAdvice(c *gin.Context) {
	var data GetAdviceData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
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
	//获取建议
	var num *int64
	advices, num,err := adminService.GetAdvices(data.PageNum, data.PageSize)
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
		var name string
		if advice.Anonymity{
			name="匿名"
		}else{
			name=student.Name
		}
		adviceResponse = append(adviceResponse, GetAdviceResponse{
			Name:        name,
			Advice:      advice.Content,
			CreatedTime: advice.CreateTime,
		})
	}
	utils.JsonSuccessResponse(c, gin.H{"data": adviceResponse, "total_page_num": math.Ceil(float64(*num)/float64(data.PageSize))})
}
