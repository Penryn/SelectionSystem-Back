package userController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/userService"
	"SelectionSystem-Back/app/utils"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

func UploadAvatar(c *gin.Context) {
	// 获取用户id
	userID, er := c.Get("ID")
	if !er {
		utils.JsonErrorResponse(c, apiException.NoThatWrong)
		return
	}
	ID, _ := userID.(int)
	// 查询用户
	user, err := userService.GetUserByID(ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 保存图片文件
	file, err := c.FormFile("avatar")
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "tempdir")
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	defer os.RemoveAll(tempDir) // 在处理完之后删除临时目录及其中的文件
	// 在临时目录中创建临时文件
	tempFile := filepath.Join(tempDir, file.Filename)
	f, err := os.Create(tempFile)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	defer f.Close()

	// 将上传的文件保存到临时文件中
	src, err := file.Open()
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	defer src.Close()

	_, err = io.Copy(f, src)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	// 判断文件的MIME类型是否为图片
	mime, err := mimetype.DetectFile(tempFile)
	if err != nil || !strings.HasPrefix(mime.String(), "image/") {
		utils.JsonErrorResponse(c, apiException.PictureError)
		return
	}

	// 保存原始图片
	filename := uuid.New().String() + ".jpg"  // 修改扩展名为.jpg
	dst := "./static/" + filename
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	// 转换图像为JPG格式并压缩
	jpgFile := filepath.Join(tempDir, "compressed.jpg")
	err = convertAndCompressImage(dst, jpgFile)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	// 替换原始文件为压缩后的JPG文件
	err = os.Rename(jpgFile, dst)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}

	url := "http://phlin.love/static/" + filename

	err = userService.UpdateAvatar(user.ID, url)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, gin.H{
		"avatar": "http://" + url,
	})
}

// 用于转换和压缩图像的函数
func convertAndCompressImage(srcPath, dstPath string) error {
    srcImg, err := imaging.Open(srcPath)
    if err != nil {
        return err
    }

    // 调整图像大小（根据需要进行调整）
    resizedImg := resize.Resize(300, 0, srcImg, resize.Lanczos3)

    // 创建新的JPG文件
    dstFile, err := os.Create(dstPath)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    // 以JPG格式保存调整大小的图像，并设置压缩质量为90
    err = jpeg.Encode(dstFile, resizedImg, &jpeg.Options{Quality: 90})
    if err != nil {
        return err
    }

    return nil
}
