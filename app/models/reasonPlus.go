package models

type ReasonPlus struct {
	Id              int64  `json:"id"`
	StudentID       string `json:"student_id"`
	TeacherID       int    `json:"teacher_id"`
	ReasonByTeacher string `json:"reason_by_teacher"`
	AdviceByTeacher string `json:"advice_by_teacher"`
	AdminID         int    `json:"admin_id"`
	ReasonByAdmin   string `json:"reason_by_admin"`
	AdviceByAdmin   string `json:"advice_by_admin"`
}
