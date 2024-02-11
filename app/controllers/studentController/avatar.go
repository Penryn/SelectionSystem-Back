package studentController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/studentService"
	"SelectionSystem-Back/app/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path"
	"strings"
)

// 上传头像
func UploadAvatar(c *gin.Context) {
	form, _ := c.MultipartForm()
	img := form.File["avatar"][0]
	imgName := img.Filename
	_ = c.SaveUploadedFile(img, "./tmp/"+imgName)
	file, _ := os.Open("./tmp/" + imgName)
	buffer := make([]byte, 512)
	_, _ = file.Read(buffer)
	contentType := http.DetectContentType(buffer)
	file.Close()
	file, _ = os.Open("./tmp/" + imgName)
	defer file.Close()
	if contentType == "image/png" {
		newTypeName := "./tmp/" + strings.TrimSuffix(img.Filename, path.Ext(path.Base(img.Filename))) + ".png"
		_ = os.Rename("./tmp/"+imgName, newTypeName)
		imgNew, err := png.Decode(file)
		if err != nil {
			fmt.Println(err)
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		out, err := os.Create("./tmp/" + strings.TrimSuffix(img.Filename, path.Ext(path.Base(img.Filename))) + ".jpg")
		if err != nil {
			fmt.Println(err)
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		defer out.Close()
		err = jpeg.Encode(out, imgNew, &jpeg.Options{Quality: 95})
		if err != nil {
			fmt.Println(err)
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		_ = os.Remove(newTypeName)
		imgName = strings.TrimSuffix(img.Filename, path.Ext(path.Base(img.Filename))) + ".jpg"
		file.Close()
		file, _ = os.Open("./tmp/" + imgName)
	}
	imgNew, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println(err)
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	fileName := uuid.NewString() + ".jpg"
	output, _ := os.Create("./avatar/" + fileName)
	defer output.Close()
	err = jpeg.Encode(output, imgNew, &jpeg.Options{Quality: 40})
	if err != nil {
		fmt.Println(err)
		return
	}
	file.Close()
	_ = os.Remove("./tmp/" + imgName)

	//获取用户身份token
	userId, er := c.Get("UserID")
	if !er {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	err = studentService.UpdateAvatar(userId.(int), "./avatar/"+fileName)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, "./avatar/"+fileName)
}
