package models

type Teacher struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`      //用户编号
	TeacherName string    `json:"teacher_name"` //教师姓名
	Section     string    `json:"section"`      //部门
	Office      string    `json:"office"`       //办公室
	Phone       string    `json:"phone"`        //电话
	Email       string    `json:"email"`        //邮箱
	Students    []Student `json:"students"`     //学生
}
