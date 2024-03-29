package studentController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/app/services/studentService"
	"SelectionSystem-Back/app/utils"
	"regexp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StudentData struct {
	Name            string `json:"name"`
	StudentID       string `json:"studentID"`
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
	TeacherName     string `json:"teacher_name"`
	TargetName      string `json:"target_name"`
	TargetAgree     int    `json:"target_agree"`
	AdminAgree      int    `json:"admin_agree"`
}

// 获取学生个人信息
func GetStudentInfo(c *gin.Context) {
	//获取用户身份token
	userId, er := c.Get("ID")
	if !er {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	user, err := studentService.GetUserByID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	var studentInfo *models.Student
	studentInfo, err = studentService.GetStudentInfoByUserID(userId.(int))
	if err != nil && err != gorm.ErrRecordNotFound {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	var targetTeacherName string
	if studentInfo.TargetID == 0 {
		targetTeacherName = "无"
	} else {
		targetTeacher, _, err := studentService.GetTeacherByTeacherID(studentInfo.TargetID)
		if err != nil && err != gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		targetTeacherName = targetTeacher.TeacherName
	}
	var ultimateTeacherName string
	if studentInfo.TeacherID == 0 {
		ultimateTeacherName = "无"
	} else {
		ultimateTeacher, _, err := studentService.GetTeacherByTeacherID(studentInfo.TeacherID)
		if err != nil && err != gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		ultimateTeacherName = ultimateTeacher.TeacherName
	}

	studentData := StudentData{
		Name:            studentInfo.Name,
		StudentID:       studentInfo.StudentID,
		Class:           studentInfo.Class,
		Phone:           studentInfo.Phone,
		PoliticalStatus: studentInfo.PoliticalStatus,
		Email:           studentInfo.Email,
		Address:         studentInfo.Address,
		Plan:            studentInfo.Plan,
		Experience:      studentInfo.Experience,
		Honor:           studentInfo.Honor,
		Interest:        studentInfo.Interest,
		Avatar:          user.Avartar,
		TeacherName:     ultimateTeacherName,
		TargetName:      targetTeacherName,
		TargetAgree:     studentInfo.TargetStatus,
		AdminAgree:      studentInfo.AdminStatus,
	}

	utils.JsonSuccessResponse(c, studentData)
}

type StudentInfoData struct {
	Name            string `json:"name" binding:"required"`
	Class           string `json:"class" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	PoliticalStatus string `json:"political_status" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Address         string `json:"address" binding:"required"`
	Plan            string `json:"plan" binding:"required"`
	Experience      string `json:"experience" binding:"required"`
	Honor           string `json:"honor" binding:"required"`
	Interest        string `json:"interest" binding:"required"`
}

// 修改个人信息
func UpdateStudentInfo(c *gin.Context) {
	var data StudentInfoData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}

	//获取用户身份token
	userId, er := c.Get("ID")
	if !er {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	var studentInfo *models.Student
	studentInfo, err = studentService.GetStudentInfoByUserID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//判断电话是否符合格式
	phone_sample := regexp.MustCompile(`^1[3|4|5|7|8][0-9]\d{8}$`)
	if !phone_sample.MatchString(data.Phone) {
		utils.JsonErrorResponse(c, apiException.PhoneError)
		return
	}
	//判断邮箱是否符合格式
	email_sample := regexp.MustCompile(`^[\w-]+(\.[\w-]+)*@[\w-]+(\.[\w-]+)+$`)
	if !email_sample.MatchString(data.Email) {
		utils.JsonErrorResponse(c, apiException.EmailError)
		return
	}
	if studentInfo.Phone != data.Phone {
		err = studentService.StudentExistByPhone(userId.(int), data.Phone)
		if err == nil {
			utils.JsonErrorResponse(c, apiException.PhoneExist)
			return
		}
	}
	if studentInfo.Email != data.Email {
		err = studentService.StudentExistByEmail(userId.(int), data.Email)
		if err == nil {
			utils.JsonErrorResponse(c, apiException.EmailExist)
			return
		}
	}

	err = studentService.UpdateStudentInfo(userId.(int), models.Student{
		Name:            data.Name,
		Class:           data.Class,
		Phone:           data.Phone,
		PoliticalStatus: data.PoliticalStatus,
		Email:           data.Email,
		Address:         data.Address,
		Plan:            data.Plan,
		Experience:      data.Experience,
		Honor:           data.Honor,
		Interest:        data.Interest,
	})
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}
