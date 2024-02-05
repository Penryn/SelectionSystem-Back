package models

type User struct {
	ID       int  `json:"id"`       //登录编号
	Username string `json:"username"` //用户名
	Password string `json:"password"` //密码
	Type     int    `json:"type"`     //用户类型  1：学生  2：教师 3：管理员
	Avartar  string `json:"avartar"`  //头像
}

