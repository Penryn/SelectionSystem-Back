package models

type Teacher struct {
	Id   int64  `json:"id"`
	TeacherID int64  `json:"teacherID"` //教师编号
	TeacherName string `json:"teacherName"` //教师姓名
	Students    []Student `json:"students"` //学生
}