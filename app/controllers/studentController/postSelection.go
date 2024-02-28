package studentController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/studentService"
	"SelectionSystem-Back/app/services/teacherService"
	"SelectionSystem-Back/app/utils"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"time"
)

type PostTeacherID struct {
	TargetID int `json:"teacher_id" binding:"required"`
}

// 提交目标导师
func PostTeacher(c *gin.Context) {
	var data PostTeacherID
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

	studentInfo, err := studentService.GetStudentInfoByUserID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	if studentInfo.Name == "未填写" {
		utils.JsonErrorResponse(c, apiException.StudentInfoWrong)
		return
	}

	targetTeacher, studentNumber, err := studentService.GetTeacherByTeacherID(data.TargetID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	teacherDDL, err := studentService.GetTeacherDDLByUserID(targetTeacher.UserID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	currentTime := time.Now()
	if currentTime.After(teacherDDL.FirstDDL) {
		utils.JsonErrorResponse(c, apiException.DDLWrong)
		return
	}

	if studentNumber >= 6 {
		utils.JsonErrorResponse(c, apiException.OverNumber)
		return
	}

	if studentInfo.TargetID != 0 && studentInfo.TargetStatus == 1 {
		originTargetTeacher, _, err := studentService.GetTeacherByTeacherID(studentInfo.TargetID)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		originTargetTeacher.StudentsNum = originTargetTeacher.StudentsNum - 1
		err = studentService.UpdateTeacher(studentInfo.TargetID, originTargetTeacher.StudentsNum)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	}

	err = studentService.UpdateTargetTeacher(userId.(int), data.TargetID, studentInfo)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	targetTeacher.StudentsNum = targetTeacher.StudentsNum + 1
	err = studentService.UpdateTeacher(data.TargetID, targetTeacher.StudentsNum)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, nil)
}

func UploadSelectionTable(c *gin.Context) {
	userId, er := c.Get("ID")
	if !er {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 查询用户
	user, err := studentService.GetUserByID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	student, err := studentService.GetStudentInfoByUserID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//查看个人信息有无填写
	if student.Name == "未填写" {
		utils.JsonErrorResponse(c, apiException.StudentInfoWrong)
		return
	}
	//查看是否在规定的期限内
	adminDDL, err := teacherService.GetAdminDDL()
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	currentTime := time.Now()
	if currentTime.After(adminDDL.SecondDDL) {
		utils.JsonErrorResponse(c, apiException.DDLWrong)
		return
	}
	//查询教师是否同意
	if student.TargetStatus != 2 {
		utils.JsonErrorResponse(c, apiException.StatusWrong)
		return
	}
	// 保存文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 判断文件的扩展名，只允许doc和docx
	ext := filepath.Ext(file.Filename)
	if ext != ".doc" && ext != ".docx" {
		utils.JsonErrorResponse(c, apiException.FileTypeInvalid)
		return
	}
	// 构建文件URL并返回
	const baseURL = "https://phlin.love/files/"
	fileName := student.Name + filepath.Ext(file.Filename)
	url := baseURL + fileName
	if _, err := os.Stat("./files"); os.IsNotExist(err) {
		err := os.Mkdir("./files", 0755)
		if err != nil {
			_ = c.AbortWithError(200, apiException.ServerError)
			return
		}
	}
	//保存
	if err := c.SaveUploadedFile(file, "./files/"+fileName); err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 更新数据库中的文件URL
	err = studentService.UpdateSelectionTable(user.ID, url)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, url)
}
