package adminController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/adminService"
	"SelectionSystem-Back/app/services/userService"
	"SelectionSystem-Back/app/utils"
	"math"

	"github.com/gin-gonic/gin"
)

type GetTableData struct {
	StudentID string `form:"student_id"`
}

type GetTableRequest struct {
	StudentID      string `json:"student_id"`
	Name           string `json:"name"`
	SelectionTable string `json:"selection_table"`
}

func GetTable(c *gin.Context) {
	var data GetTableData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//获取用户id
	userID, er := c.Get("ID")
	if !er {
		utils.JsonErrorResponse(c, apiException.NoThatWrong)
		return
	}
	ID, _ := userID.(int)
	//查询用户
	user, err := userService.GetUserByID(ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//鉴权
	if user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//查询学生
	student, err := userService.GetStudentByStudentID(data.StudentID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, gin.H{"student_id": student.StudentID, "name": student.Name, "selection_table": student.SelectionTable})
}

type CheckTableData struct {
	StudentsID []string `json:"students_id" validate:"required"`
	ReasonID  int      `json:"reason_id"`
	Check      int      `json:"check" validate:"oneof=1 2"` // 1:同意 2:拒绝
}

func CheckTable(c *gin.Context) {
	var data CheckTableData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//获取用户id
	userID, er := c.Get("ID")
	if !er {
		utils.JsonErrorResponse(c, apiException.NoThatWrong)
		return
	}
	ID, _ := userID.(int)
	//查询用户
	user, err := userService.GetUserByID(ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//鉴权
	if user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	if data.Check==2 && data.ReasonID==0{
		utils.JsonErrorResponse(c, apiException.ReasonError)
		return
	}
	//批量处理学生
	if len(data.StudentsID) >6 {
		utils.JsonErrorResponse(c, apiException.MoreThanSix)
		return
	}
	for _, studentID := range data.StudentsID {

		//查询学生
		student, err := userService.GetStudentByStudentID(studentID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		if student.TargetStatus != 2 {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		if student.AdminStatus != 1 {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		err = adminService.CheckTable(student.StudentID, student.TargetID, data.Check)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		if data.Check==2{
			reason ,err:= userService.GetReasonsByReasonID(data.ReasonID)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
			err = userService.SendConversation(user.ID, student.UserID, "您的双向申请表申请被拒绝，原因："+reason.ReasonContent)
			if err != nil {
				utils.JsonErrorResponse(c, apiException.ServerError)
				return
			}
		}
	}
	utils.JsonSuccessResponse(c, nil)
}


type GetPostData struct {
	Check int `form:"check" validate:"oneof=1 2"` // 1:待处理 2:已处理
	Name string `form:"name"`
	StudentID string `form:"student_id"`
}

type GetPostResponse struct {
	StudentID string `json:"student_id"`
	Name      string `json:"name"`
	Status   int    `json:"status"`
}

func GetPost(c *gin.Context) {
	var data GetPostData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//获取用户id
	userID, er := c.Get("ID")
	if !er {
		utils.JsonErrorResponse(c, apiException.NoThatWrong)
		return
	}
	ID, _ := userID.(int)
	//查询用户
	user, err := userService.GetUserByID(ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//鉴权
	if user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	students, err := adminService.GetCheckStudents(data.Check,data.Name,data.StudentID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	result:=make([]GetPostResponse,0)
	for i := 0; i < len(students); i++ {
		if students[i].AdminStatus == 1 {
			result = append(result, GetPostResponse{StudentID: students[i].StudentID, Name: students[i].Name, Status: 0})
		} else if students[i].AdminStatus == 2 {
			result = append(result, GetPostResponse{StudentID: students[i].StudentID, Name: students[i].Name, Status: 1})
		}else if students[i].AdminStatus == 3 {
			result = append(result, GetPostResponse{StudentID: students[i].StudentID, Name: students[i].Name, Status: 2})
		}

	}
	utils.JsonSuccessResponse(c, gin.H{"data": result})
}

type DisassociateData struct {
	StudentID string `json:"student_id"`
}

func Disassociate(c *gin.Context) {
	var data DisassociateData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//获取用户id
	userID, er := c.Get("ID")
	if !er {
		utils.JsonErrorResponse(c, apiException.NoThatWrong)
		return
	}
	ID, _ := userID.(int)
	//查询用户
	user, err := userService.GetUserByID(ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//鉴权
	if user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//查询学生
	student, err := userService.GetStudentByStudentID(data.StudentID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	if student.AdminStatus != 2 {
		utils.JsonErrorResponse(c, apiException.AdminStatusError)
		return
	}
	err = adminService.Disassociate(student.StudentID, student.TargetID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, nil)
}

type GetTeacherWithStudentsData struct {
	PageNum  int `form:"page_num" validate:"required"`
	PageSize int `form:"page_size" validate:"required"`
}

type Student struct {
	StudentName string `json:"student_name"`
	StudentID   string `json:"student_id"`
}

type Request struct {
	TeacherID    int       `json:"teacher_id"`
	TeacherName  string    `json:"teacher_name"`
	Students     []Student `json:"students"`
}

func GetTeacherWithStudents(c *gin.Context) {
	var data GetTeacherWithStudentsData
	err := c.ShouldBindQuery(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//获取用户id
	userID, er := c.Get("ID")
	if !er {
		utils.JsonErrorResponse(c, apiException.NoThatWrong)
		return
	}
	ID, _ := userID.(int)
	//查询用户
	user, err := userService.GetUserByID(ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//鉴权
	if user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	var num *int64
	teachers, num,err := adminService.GetTeachers(data.PageNum, data.PageSize)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	result := make([]Request, 0)
	for i := 0; i < len(teachers); i++ {
		var students []Student
		for j := 0; j < len(teachers[i].Students); j++ {
			students = append(students, Student{StudentName: teachers[i].Students[j].Name, StudentID: teachers[i].Students[j].StudentID})
		}
		result = append(result, Request{TeacherID: teachers[i].ID, TeacherName: teachers[i].TeacherName, Students: students})
	}
	utils.JsonSuccessResponse(c, gin.H{"data": result, "total_page_num":math.Ceil(float64(*num)/float64(data.PageSize)) })
}