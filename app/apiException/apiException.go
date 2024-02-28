package apiException

import "net/http"

type Error struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
}

var (
	ServerError      = NewError(http.StatusInternalServerError, 200500, "系统异常，请稍后重试!")
	OpenIDError      = NewError(http.StatusInternalServerError, 200500, "系统异常，请稍后重试!")
	ParamError       = NewError(http.StatusInternalServerError, 200501, "参数错误")
	NoThatWrong      = NewError(http.StatusInternalServerError, 200502, "账号、密码或身份错误")
	UserNotFind      = NewError(http.StatusNotFound, 200503, "用户不存在")
	PictureError     = NewError(http.StatusInternalServerError, 200504, "仅允许上传图片文件")
	MoreThanSix      = NewError(http.StatusInternalServerError, 200505, "一次最多处理6个学生")
	TimeSetError     = NewError(http.StatusInternalServerError, 200506, "时间设置错误")
	EmailExist       = NewError(http.StatusInternalServerError, 200507, "邮箱已存在")
	PhoneExist       = NewError(http.StatusInternalServerError, 200508, "电话已存在")
	DDLWrong         = NewError(http.StatusInternalServerError, 200509, "已超过最晚期限")
	OverNumber       = NewError(http.StatusInternalServerError, 200510, "超过限制人数")
	FileTypeInvalid  = NewError(http.StatusInternalServerError, 200511, "只允许上传docx/doc文件")
	StudentInfoWrong = NewError(http.StatusInternalServerError, 200512, "未填写个人信息")
	StatusWrong      = NewError(http.StatusInternalServerError, 200513, "教师暂未同意你的审批请求")
	StudentNotFound  = NewError(http.StatusInternalServerError, 200514, "您没有权限对该学生的请求进行审批")
	StudentWrong     = NewError(http.StatusInternalServerError, 200515, "该学生并非您的学生")
	NotFound         = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
	Unknown          = NewError(http.StatusInternalServerError, 300500, "系统异常，请稍后重试!")
)

func OtherError(message string) *Error {
	return NewError(http.StatusForbidden, 100403, message)
}

func (e *Error) Error() string {
	return e.Msg
}

func NewError(statusCode, Code int, msg string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
	}
}
