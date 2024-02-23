package teacherController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/teacherService"
	"SelectionSystem-Back/app/utils"
	"github.com/gin-gonic/gin"
	"strconv"
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

// 获取学生列表
func GetStudentList(c *gin.Context) {
	check := c.Query("check")
	checkStudentList, err := strconv.Atoi(check)
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

	teacher, _, err := teacherService.GetTeacherByUserID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	studentList, err := teacherService.StudentList(teacher.ID, checkStudentList)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	var responseStudentList = make([]Student, 0)
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

type UltimateStudent struct {
	Name      string `json:"name" binding:"required"`
	StudentID string `json:"student_id" binding:"required"`
}

// 获取最终学生
func GetUltimateStudentList(c *gin.Context) {
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

	students, num, err := teacherService.GetStudentsByUserID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	var responseStudentList = make([]UltimateStudent, 0)
	for _, student := range students {
		response := UltimateStudent{
			Name:      student.Name,
			StudentID: student.StudentID,
		}
		responseStudentList = append(responseStudentList, response)
	}
	utils.JsonSuccessResponse(c, gin.H{
		"student_num": num,
		"data":        responseStudentList,
	})
}

type StudentData struct {
	Name            string `json:"name"`
	StudentID       string `json:"student_id"`
	Class           string `json:"class"`
	Phone           string `json:"phone"`
	PoliticalStatus string `json:"political_status"`
	Email           string `json:"email"`
	Address         string `json:"address"`
	Plan            string `json:"plan"`
	Experience      string `json:"experience"`
	Honor           string `json:"honor"`
	Interest        string `json:"interest"`
	Avatar          string `json:"avartar"`
}

// 获取学生信息
func GetStudentInfo(c *gin.Context) {
	studentId := c.Query("student_id")
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

	student, err := teacherService.GetStudentInfoByStudentID(studentId)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	studentInfo, err := teacherService.GetUserByID(student.UserID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	response := StudentData{
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
	}

	utils.JsonSuccessResponse(c, response)
}

type MessagedStudent struct {
	UserID int    `json:"user_id" binding:"required"`
	Name   string `json:"name" binding:"Required"`
}

// 获取私聊过自己的学生列表
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

	if user.Type != 2 && user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	conversations, err := teacherService.GetMessagedStudentListByUserID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	var responseStudentList = make([]MessagedStudent, 0)
	for _, conversation := range conversations {
		studentInfo, err := teacherService.GetStudentInfoByUserID(conversation.UserAID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		response := MessagedStudent{
			UserID: conversation.UserAID,
			Name:   studentInfo.Name,
		}
		responseStudentList = append(responseStudentList, response)
	}

	utils.JsonSuccessResponse(c, responseStudentList)
}
