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

func AdviceGet(c *gin.Context) {
	userId, er := c.Get("ID")
	if !er {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	var advice []models.Advice
	advice, err := studentService.GetAdvice(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, advice)
}
