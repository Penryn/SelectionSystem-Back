package apiException

import "net/http"

type Error struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
}

var (
	ServerError              = NewError(http.StatusInternalServerError, 200500, "系统异常，请稍后重试!")
	OpenIDError              = NewError(http.StatusInternalServerError, 200500, "系统异常，请稍后重试!")
	ParamError               = NewError(http.StatusInternalServerError, 200501, "参数错误")
	NoThatWrong              = NewError(http.StatusInternalServerError, 200502, "账号、密码或身份错误")
	UserNotFind              = NewError(http.StatusNotFound, 200503, "用户不存在")
	PictureError             = NewError(http.StatusInternalServerError, 200504, "仅允许上传图片文件")
	MoreThanSix              = NewError(http.StatusInternalServerError, 200505, "一次最多处理6个学生")
	TimeSetError             = NewError(http.StatusInternalServerError, 200506, "时间设置错误")
	EmailExist               = NewError(http.StatusInternalServerError, 200507, "邮箱已存在")
	PhoneExist               = NewError(http.StatusInternalServerError, 200508, "电话已存在")
	DDLWrong                 = NewError(http.StatusInternalServerError, 200509, "已超过最晚期限")
	OverNumber               = NewError(http.StatusInternalServerError, 200510, "超过限制人数")
	FileTypeInvalid          = NewError(http.StatusInternalServerError, 200511, "只允许上传docx/doc文件")
	StudentInfoWrong         = NewError(http.StatusInternalServerError, 200512, "未填写个人信息")
	StatusWrong              = NewError(http.StatusInternalServerError, 200513, "教师暂未同意你的审批请求")
	StudentNotFound          = NewError(http.StatusInternalServerError, 200514, "您没有权限对该学生的请求进行审批")
	StudentWrong             = NewError(http.StatusInternalServerError, 200515, "该学生并非您的学生")
	ReasonError              = NewError(http.StatusInternalServerError, 200516, "拒绝时必须填写拒绝理由")
	AdminStatusError         = NewError(http.StatusInternalServerError, 200517, "该用户并未与老师双向选择")
	ReasonExist              = NewError(http.StatusInternalServerError, 200518, "原因名称已存在")
	MessageError             = NewError(http.StatusInternalServerError, 200519, "消息不能为空")
	StudentExistError        = NewError(http.StatusInternalServerError, 200520, "该学号不存在")
	ReasonNameOrContentEmpty = NewError(http.StatusInternalServerError, 200521, "原因名称或原因内容不能为空")
	ReasonExistError         = NewError(http.StatusInternalServerError, 200522, "该原因不存在")
	TeacherNotFound          = NewError(http.StatusInternalServerError, 200522, "该老师不存在")
	AdminPostError           = NewError(http.StatusInternalServerError, 200523, "该学生并未提交申请或已被审批")
	TeacherPostError         = NewError(http.StatusInternalServerError, 200524, "该学生未被教师审批通过")
	AdminError               = NewError(http.StatusInternalServerError, 200525, "该学生已被管理员审批通过，若需要请到最终学生处解除关系")
	NotFound                 = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
	Unknown                  = NewError(http.StatusInternalServerError, 300500, "系统异常，请稍后重试!")
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
