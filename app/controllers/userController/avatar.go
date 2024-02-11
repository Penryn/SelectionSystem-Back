package userController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/userService"
	"SelectionSystem-Back/app/utils"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadAvatar(c *gin.Context) {
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
		fmt.Println(1)
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//保存图片文件
	file, err := c.FormFile("avatar")
	if err != nil {
		fmt.Println(2)
		utils.JsonErrorResponse(c,apiException.ServerError)
		return
	}
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "tempdir")
	if err != nil {
		fmt.Println(3)
		utils.JsonErrorResponse(c,apiException.ServerError)
		return
	}
	defer os.RemoveAll(tempDir) // 在处理完之后删除临时目录及其中的文件
	// 在临时目录中创建临时文件
	tempFile := filepath.Join(tempDir, file.Filename)
	f, err := os.Create(tempFile)
	if err != nil {
		fmt.Println(4)
		utils.JsonErrorResponse(c,apiException.ServerError)
		return
	}
	defer f.Close()

	// 将上传的文件保存到临时文件中
	src, err := file.Open()
	if err != nil {
		fmt.Println(5)
		utils.JsonErrorResponse(c,apiException.ServerError)
		return
	}
	defer src.Close()

	_, err = io.Copy(f, src)
	if err != nil {
		fmt.Println(6)
		utils.JsonErrorResponse(c,apiException.ServerError)
		return
	}
	// 判断文件的MIME类型是否为图片
	mime, err := mimetype.DetectFile(tempFile)
	if err != nil || !strings.HasPrefix(mime.String(), "image/") {
		fmt.Println(7)
		utils.JsonErrorResponse(c, apiException.PictureError)
		return
	}
	//继续保存图片
	filename := uuid.New().String() + filepath.Ext(file.Filename)
	dst := "./static/" + filename
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		fmt.Println(8)
		utils.JsonErrorResponse(c,apiException.ServerError)
		return
	}
	url := "47.115.209.120:8080" + "/static/" + filename

	err = userService.UpdateAvatar(user.ID, url)
	if err != nil {
		fmt.Println(9)
		utils.JsonErrorResponse(c,apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, gin.H{
		"avatar": "http://" + url,
	})
}