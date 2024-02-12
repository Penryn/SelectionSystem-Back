package studentController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/app/services/studentService"
	"SelectionSystem-Back/app/utils"
	"github.com/gin-gonic/gin"
)

type AdviceData struct {
	Advice    string `json:"advice" binding:"required"`
	Anonymity bool   `json:"anonymity" binding:"required"`
}

func AdvicePost(c *gin.Context) {
	var data AdviceData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}

	userId, er := c.Get("UserID")
	if !er {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	err = studentService.CreateAdvice(models.Advice{
		UserID:    userId.(int),
		Content:   data.Advice,
		Anonymity: data.Anonymity,
	})
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

func AdviceGet(c *gin.Context) {
	userId, er := c.Get("UserID")
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
