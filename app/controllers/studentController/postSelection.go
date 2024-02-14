package studentController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/studentService"
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

	studentInfo, err := studentService.GetStudentInfoByUserID(userId.(int))
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	err = studentService.UpdateTargetTeacher(userId.(int), data.TargetID, studentInfo)
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
	const baseURL = "http://47.115.209.120:8080/files/"
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
