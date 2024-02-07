package userController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/app/services/userService"
	"SelectionSystem-Back/app/utils"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Type     int    `json:"type"`
}

func Login(c *gin.Context) {
	//获取参数
	var data UserData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
		return
	}
	var user models.User
	var msg string
	//判断用户类型
	//学生
	if data.Type == 1 {
		// 正则表达式匹配12位数字
		pattern := "^[0-9A-Za-z]{12}$"
		_, err := regexp.MatchString(pattern, data.Username)
		if err != nil {
			utils.JsonErrorResponse(c, apiException.NoThatWrong)
			return
		}
		//如果username第一位数字是2说明是2023级以前的，否则是2023级以后的
		if data.Username[0] == '2' {
			msg = data.Username[:4] + "级学生"
		} else if data.Username[0] == '3' {
			msg = data.Username[2:6] + "级学生"
		} else {
			utils.JsonErrorResponse(c, apiException.NoThatWrong)
			return
		}
		//查找用户是否存在
		user, err = userService.GetUserByUsername(data.Username)
		if err != nil && err == gorm.ErrRecordNotFound {
			//新用户
			// 判断密码是否符合规则
			if data.Password == "zjut"+data.Username[len(data.Username)-6:] {
				// 创建新用户
				err := userService.CreateUser(models.User{
					Username: data.Username,
					Password: data.Password,
					Type:     data.Type,
					Avartar:  userService.GetAvartar(),
				})
				if err != nil {
					utils.JsonErrorResponse(c, apiException.ServerError)
					return
				}
				// 返回新用户信息
				user, err = userService.GetUserByUsername(data.Username)
				if err != nil {
					utils.JsonErrorResponse(c, apiException.ServerError)
					return
				}
			} else {
				utils.JsonErrorResponse(c, apiException.NoThatWrong)
				return
			}
		} else if err == nil {
			//用户存在
			//判断密码是否正确
			if user.Password != data.Password {
				utils.JsonErrorResponse(c, apiException.NoThatWrong)
				return
			}
		} else {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
		//教师和管理员
	} else if data.Type == 2 || data.Type == 3 {
		if data.Type == 2 {
			msg = "教师"
		} else {
			msg = "管理员"
		}
		//查找用户是否存在
		user, err = userService.GetUserByUsername(data.Username)
		if err != nil && err == gorm.ErrRecordNotFound {
			utils.JsonErrorResponse(c, apiException.UserNotFind)
			return
		} else if err == nil {
			if user.Password != data.Password {
				utils.JsonErrorResponse(c, apiException.NoThatWrong)
				return
			}
		} else {
			utils.JsonErrorResponse(c, apiException.ServerError)
			return
		}
	} else {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//获取jwt
	token, err := utils.GenToken(user.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	type loginresult struct {
		Token   string `json:"token"`
		Msg     string `json:"msg"`
		Avartar string `json:"avartar"`
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": loginresult{
			Token:   token,
			Msg:     msg,
			Avartar: user.Avartar},
	})
}
