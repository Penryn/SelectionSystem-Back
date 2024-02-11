package studentController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/app/services/studentService"
	"SelectionSystem-Back/app/utils"
	"github.com/gin-gonic/gin"
	"time"
)

type PageData struct {
	PageNum  int `form:"page_num" binding:"required"`
	PageSize int `form:"page_size" binding:"required"`
}

// 获取教师列表
func GetTeacherList(c *gin.Context) {
	var data PageData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}

	teacherList, err := studentService.GetTeacherList()
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	//获取用户身份token
	userId, er := c.Get("UserID")
	if !er {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	adminDDL, err := studentService.GetAdminDDL(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	currentTime := time.Now()
	if currentTime.After(adminDDL.FirstDDL) {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	var responseTeacherList []models.Teacher
	for _, teacher := range teacherList {
		studentCount, err := studentService.CheckTeacherList(teacher)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		if studentCount >= 6 {
			continue
		}
		response := models.Teacher{
			ID:          teacher.ID,
			TeacherName: teacher.TeacherName,
			Section:     teacher.Section,
			Office:      teacher.Office,
			Phone:       teacher.Phone,
			Email:       teacher.Email,
		}
		responseTeacherList = append(responseTeacherList, response)
	}

	pageTeacherList, totalPageNum := studentService.PageTeacherList(responseTeacherList, data.PageNum, data.PageSize)

	utils.JsonSuccessResponse(c, gin.H{
		"total_page_num": totalPageNum,
		"data":           pageTeacherList,
	})
}
