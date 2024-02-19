package studentController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/app/services/studentService"
	"SelectionSystem-Back/app/utils"
	"github.com/gin-gonic/gin"
	"time"
)

type AdviceData struct {
	Advice    string `json:"advice" binding:"required"`
	Anonymity *bool  `json:"anonymity"`
}

func AdvicePost(c *gin.Context) {
	var data AdviceData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}

	userId, er := c.Get("ID")
	if !er {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	student, err := studentService.GetStudentInfoByUserID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	if student.Name == "未填写" {
		utils.JsonErrorResponse(c, apiException.StudentInfoWrong)
		return
	}

	if data.Anonymity == nil {
		anonymity := false
		data.Anonymity = &anonymity
	}

	err = studentService.CreateAdvice(models.Advice{
		UserID:     userId.(int),
		Content:    data.Advice,
		Anonymity:  *data.Anonymity,
		CreateTime: time.Now(),
	})
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

type Advice struct {
	Content    string `json:"content" binding:"required"`
	Anonymity  bool   `json:"anonymity" binding:"required"`
	CreateTime string `json:"create_time" binding:"required"`
}

func AdviceGet(c *gin.Context) {
	userId, er := c.Get("ID")
	if !er {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	advice, err := studentService.GetAdvice(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	var responseAdvice = make([]Advice, 0)
	for _, adviceInfo := range advice {
		formattedTime := adviceInfo.CreateTime.Format("2006-01-02T15:04:05Z")
		response := Advice{
			Content:    adviceInfo.Content,
			Anonymity:  adviceInfo.Anonymity,
			CreateTime: formattedTime,
		}
		responseAdvice = append(responseAdvice, response)
	}

	utils.JsonSuccessResponse(c, responseAdvice)
}
