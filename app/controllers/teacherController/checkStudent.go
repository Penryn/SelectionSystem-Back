package teacherController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/teacherService"
	"SelectionSystem-Back/app/utils"
	"github.com/gin-gonic/gin"
	"time"
)

// 教师审批
type CheckData struct {
	StudentsID []string `json:"students_id" binding:"required"`
	Check      int      `json:"check" binding:"required"` // 2:同意 3:拒绝
}

func CheckByTeacher(c *gin.Context) {
	var data CheckData
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

	teacher, _, err := teacherService.GetTeacherByUserID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	if data.Check == 1 {
		studentsNum, err := teacherService.GetStudentsNumByTarget(teacher.ID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		if int64(len(data.StudentsID)) > 6-studentsNum {
			utils.JsonErrorResponse(c, apiException.OverNumber)
			return
		}
	}

	if data.Check == 1 {
		for _, studentId := range data.StudentsID {
			studentInfo, err := teacherService.GetStudentInfoByStudentID(studentId)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
			studentInfo.TargetStatus = 2
			err = teacherService.UpdateStudentInfo(studentId, studentInfo)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
			if studentInfo.AdminStatus == 2 {
				err = teacherService.StudentJoinTeacher(studentId, studentInfo.TargetID)
				if err != nil {
					utils.JsonErrorResponse(c, apiException.ServerError)
					return
				}
			}
		}
	} else if data.Check == 2 {
		for _, studentId := range data.StudentsID {
			studentInfo, err := teacherService.GetStudentInfoByStudentID(studentId)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
			studentInfo.TargetStatus = 3
			teacher.StudentsNum = teacher.StudentsNum - 1
			err = teacherService.UpdateTeacher(userId.(int), teacher.StudentsNum)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
			err = teacherService.UpdateStudentInfo(studentId, studentInfo)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
		}
	} else {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
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
		utils.JsonErrorResponse(c, apiException.DDLWrong)
		return
	}

	studentInfo, err := teacherService.GetStudentInfoByStudentID(data.StudentID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	if studentInfo.TargetStatus != 2 || studentInfo.AdminStatus != 2 {
		utils.JsonErrorResponse(c, apiException.StatusWrong)
		return
	}

	err = teacherService.Disassociate(data.StudentID, studentInfo.TeacherID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	teacher, _, err := teacherService.GetTeacherByUserID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	teacher.StudentsNum = teacher.StudentsNum - 1
	err = teacherService.UpdateTeacher(userId.(int), teacher.StudentsNum)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

// 教师撤回审批
type WithdrawData struct {
	StudentsID []string `json:"students_id" binding:"required"`
}

func WithdrawApproval(c *gin.Context) {
	var data WithdrawData
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
		utils.JsonErrorResponse(c, apiException.DDLWrong)
		return
	}

	for _, studentId := range data.StudentsID {
		studentInfo, err := teacherService.GetStudentInfoByStudentID(studentId)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		teacher, _, err := teacherService.GetTeacherByID(studentInfo.TargetID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		if studentInfo.TargetStatus == 2 {
			studentInfo.TargetStatus = 1
		} else if studentInfo.TargetStatus == 3 {
			teacher.StudentsNum = teacher.StudentsNum + 1
			err = teacherService.UpdateTeacher(userId.(int), teacher.StudentsNum)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
			studentInfo.TargetStatus = 1
		} else {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		err = teacherService.UpdateStudentInfoTargetStatus(studentId, studentInfo.TargetStatus)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	}

	utils.JsonSuccessResponse(c, nil)
}
