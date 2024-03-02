package userController

import (
	"SelectionSystem-Back/app/apiException"
	"SelectionSystem-Back/app/services/userService"
	"SelectionSystem-Back/app/utils"

	"github.com/gin-gonic/gin"
)

type CreateReasonData struct {
	ReasonName    string `json:"reason_name"`    //原因名称
	ReasonContent string `json:"reason_content"` //原因内容
}

func CreateReason(c *gin.Context) {
	var data CreateReasonData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
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
	if user.Type != 2 && user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//不能为空
	if data.ReasonName == "" || data.ReasonContent == "" {
		utils.JsonErrorResponse(c, apiException.ReasonNameOrContentEmpty)
		return
	}
	//查找原因是否存在
	_, err = userService.GetReasonByName(data.ReasonName)
	if err == nil {
		utils.JsonErrorResponse(c, apiException.ReasonExist)
		return
	}
	//创建原因
	err = userService.CreateReason(user.ID, data.ReasonName, data.ReasonContent)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, "创建成功")
}

type UpdateReasonData struct {
	ReasonID      int    `json:"reason_id"`
	ReasonName    string `json:"reason_name"`    //原因名称
	ReasonContent string `json:"reason_content"` //原因内容

}

func UpdateReason(c *gin.Context) {
	var data UpdateReasonData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
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
	if user.Type != 2 && user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//查询原因是否存在
	_, err = userService.GetReasonByID(data.ReasonID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//更新原因
	err = userService.UpdateReason(user.ID, data.ReasonID, data.ReasonName, data.ReasonContent)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, "修改成功")
}

type DeleteReasonData struct {
	ReasonID int `json:"reason_id"`
}

func DeleteReason(c *gin.Context) {
	var data DeleteReasonData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ParamError)
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
	if user.Type != 2 && user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//删除原因
	err = userService.DeleteReason(user.ID, data.ReasonID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, "删除成功")
}

type GetReasonResponse struct {
	ID            int    `json:"id"`
	ReasonName    string `json:"reason_name"`    //原因名称
	ReasonContent string `json:"reason_content"` //原因内容
}

func GetReasons(c *gin.Context) {
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
	if user.Type != 2 && user.Type != 3 {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	//查询原因
	reasons, err := userService.GetReasons(user.ID)
	if err != nil {
		utils.JsonErrorResponse(c, apiException.ServerError)
		return
	}
	var response []GetReasonResponse
	for _, reason := range reasons {
		response = append(response, GetReasonResponse{ID: reason.ID, ReasonName: reason.ReasonName, ReasonContent: reason.ReasonContent})
	}
	utils.JsonSuccessResponse(c, response)
}
