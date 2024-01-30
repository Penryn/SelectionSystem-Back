package models

type CheckPost struct {
	ID            int    `json:"id"`              //审批id
	StudentID     string `json:"student_id"`      //学号
	Name          string `json:"name"`            //姓名
	Class         string `json:"class"`           //专业班级
	Phone         string `json:"phone"`           //电话
	TeacherID     int    `json:"teacher_id"`      //教师id
	TeacherName   string `json:"teacher_name"`    //教师名字
	AdminID       int    `json:"admin_id"`        //管理员id
	AdminName     string `json:"admin_name"`      //管理员名字
	File1         string `json:"file_1"`          //第一轮双向选择表
	File2         string `json:"file_2"`          //第二轮双向选择表
	Status1       int    `json:"status_1"`        //教师审批
	Status2       int    `json:"status_2"`        //管理员审批
	TimeByTeacher string `json:"time_by_teacher"` //教师设置第一轮审批时间
	TimeByAdmin1  string `json:"time_by_admin_1"` //管理员设置第一轮审批时间
	TimeByAdmin2  string `json:"time_by_admin_2"` //管理员设置第二轮审批时间
}
