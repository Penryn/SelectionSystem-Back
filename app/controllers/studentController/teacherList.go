package studentController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/studentService"
	"SelectionSystem-Back/app/utils"
	"github.com/gin-gonic/gin"
	"math"
)

type PageData struct {
	PageNum  int    `form:"page_num" binding:"required"`
	PageSize int    `form:"page_size" binding:"required"`
	Name     string `form:"name"`
}

type Teacher struct {
	ID          int    `json:"teacher_id" binding:"required"`
	Name        string `json:"teacherName" binding:"required"`
	Section     string `json:"section" binding:"required"`
	Office      string `json:"office" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	Email       string `json:"email" binding:"required"`
	StudentsNum int    `json:"students_num" binding:"required"`
	TeacherDDL  string `json:"teacher_ddl" binding:"required"`
}

// 获取教师列表
func GetTeacherList(c *gin.Context) {
	var data PageData
	err := c.ShouldBindQuery(&data)
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

	flag := studentService.CheckStudentInfo(student)
	if !flag {
		utils.JsonErrorResponse(c, apiException.StudentInfoWrong)
		return
	}
	
	teacherList, err := studentService.GetTeacherList(data.PageNum, data.PageSize, data.Name)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.TeacherNotFound)
		return
	}

	var totalPageNum *int64
	totalPageNum, err = studentService.GetTotalPageNum(data.Name)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	var responseTeacherList []Teacher
	for _, teacher := range teacherList {
		_, studentsNum, err := studentService.GetTeacherByTeacherID(teacher.ID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		teacherDDL, err := studentService.GetTeacherDDLByUserID(teacher.UserID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		formattedTime := teacherDDL.FirstDDL.Format("2006-01-02T15:04:05Z")
		response := Teacher{
			ID:          teacher.ID,
			Name:        teacher.TeacherName,
			Section:     teacher.Section,
			Office:      teacher.Office,
			Phone:       teacher.Phone,
			Email:       teacher.Email,
			StudentsNum: studentsNum,
			TeacherDDL:  formattedTime,
		}
		responseTeacherList = append(responseTeacherList, response)
	}

	utils.JsonSuccessResponse(c, gin.H{
		"total_page_num": math.Ceil(float64(*totalPageNum) / float64(data.PageSize)),
		"data":           responseTeacherList,
	})
}
