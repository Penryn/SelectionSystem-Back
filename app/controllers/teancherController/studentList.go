package teancherController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/teacherService"
	"SelectionSystem-Back/app/utils"
	"github.com/gin-gonic/gin"
)

type Student struct {
	Name            string `json:"name" binding:"required"`
	StudentID       string `json:"student_id" binding:"required"`
	Class           string `json:"class" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	PoliticalStatus string `json:"political_status" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Address         string `json:"address" binding:"required"`
	Plan            string `json:"plan" binding:"required"`
	Experience      string `json:"experience" binding:"required"`
	Honor           string `json:"honor" binding:"required"`
	Interest        string `json:"interest" binding:"required"`
	Avatar          string `json:"avartar" binding:"required"`
	TargetStatus    int    `json:"target_agree" binding:"required"`
	AdminStatus     int    `json:"admin_agree" binding:"required"`
}

// 获取未审批的学生列表
func GetStudentList(c *gin.Context) {
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

	if user.Type != 2 && user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	studentList, err := teacherService.StudentList(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	var responseStudentList []Student
	for _, student := range studentList {
		studentInfo, err := teacherService.GetUserByID(student.UserID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}

		response := Student{
			StudentID:       student.StudentID,
			Name:            student.Name,
			Class:           student.Class,
			Phone:           student.Phone,
			PoliticalStatus: student.PoliticalStatus,
			Email:           student.Email,
			Address:         student.Address,
			Plan:            student.Plan,
			Experience:      student.Experience,
			Honor:           student.Honor,
			Interest:        student.Interest,
			Avatar:          studentInfo.Avartar,
			TargetStatus:    student.TargetStatus,
			AdminStatus:     student.AdminStatus,
		}
		responseStudentList = append(responseStudentList, response)
	}

	utils.JsonSuccessResponse(c, responseStudentList)
}

// 获取已审批的学生列表
func GetCheckStudentList(c *gin.Context) {
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

	if user.Type != 2 && user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	studentList, err := teacherService.StudentCheckList(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	var responseStudentList []Student
	for _, student := range studentList {
		studentInfo, err := teacherService.GetUserByID(student.UserID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}

		response := Student{
			StudentID:       student.StudentID,
			Name:            student.Name,
			Class:           student.Class,
			Phone:           student.Phone,
			PoliticalStatus: student.PoliticalStatus,
			Email:           student.Email,
			Address:         student.Address,
			Plan:            student.Plan,
			Experience:      student.Experience,
			Honor:           student.Honor,
			Interest:        student.Interest,
			Avatar:          studentInfo.Avartar,
			TargetStatus:    student.TargetStatus,
			AdminStatus:     student.AdminStatus,
		}
		responseStudentList = append(responseStudentList, response)
	}

	utils.JsonSuccessResponse(c, responseStudentList)
}
