package teancherController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/app/services/teacherService"
	"SelectionSystem-Back/app/utils"
	"github.com/gin-gonic/gin"
	"time"
)

// 教师审批
type CheckData struct {
	CheckFirstPost int `json:"check_firstpost" binding:"required"`
	ID             int `json:"id" binding:"required"`
}

type CheckStudent struct {
	Checks []CheckData `json:"checks" binding:"required"`
}

func CheckByTeacher(c *gin.Context) {
	var data CheckStudent
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

	user, err := teacherService.GetUserByID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	if user.Type != 2 && user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	_, studentNumber, err := teacherService.GetTeacherByUserID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	var count = 0
	for range data.Checks {
		count++
		if count > 6-studentNumber {
			utils.JsonErrorResponse(c, apiException.OverNumber)
			return
		}
	}

	for _, check := range data.Checks {
		studentInfo, err := teacherService.GetStudentInfoByID(check.ID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}

		if check.CheckFirstPost == 1 {
			studentInfo.TargetStatus = 1
		} else if check.CheckFirstPost == 2 {
			studentInfo.TargetStatus = 2
		} else {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		err = teacherService.UpdateStudentInfo(check.ID, studentInfo)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	}

	utils.JsonSuccessResponse(c, nil)
}

type CancelStudentData struct {
	StudentID string `json:"student_id" binding:"required"`
}

func CancelStudent(c *gin.Context) {
	var data CancelStudentData
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

	user, err := teacherService.GetUserByID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	if user.Type != 2 && user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	adminDDL, err := teacherService.GetAdminDDL()
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	currentTime := time.Now()
	if currentTime.After(adminDDL.FirstDDL) {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	studentInfo, err := teacherService.GetStudentInfoByStudentID(data.StudentID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	if studentInfo.TargetStatus != 1 || studentInfo.AdminStatus != 1 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	teacher, err := teacherService.GetTeacherByTeacherID(studentInfo.TeacherID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	students, err := teacherService.GetStudentList(studentInfo.TeacherID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	updatedStudents := make([]models.Student, 0)
	for i := range students {
		if students[i].StudentID != data.StudentID {
			teacherService.AseDecryptStudentInfo(&students[i])
			updatedStudents = append(updatedStudents, students[i])
		}
	}
	teacher.Students = updatedStudents
	teacher.StudentsNum = len(updatedStudents)
	studentInfo.TargetStatus = 2
	//var teacherID *int = nil
	//studentInfo.TeacherID = teacherID
	err = teacherService.UpdateStudentInfoByStudentID(data.StudentID, studentInfo)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	err = teacherService.UpdateTeacher(teacher)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}
