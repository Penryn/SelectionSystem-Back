package models

type Student struct {
	ID              int     `json:"id"`               //用户id
	UserID          int     `json:"user_id"`          //用户编号
	StudentID       string  `json:"student_id"`       //学号
	Name            string  `json:"name"`             //姓名
	Email           string  `json:"email"`            //邮箱
	Class           string  `json:"class"`            //专业班级
	Phone           string  `json:"phone"`            //电话
	PoliticalStatus string  `json:"political_status"` //政治面貌
	Address         string  `json:"address"`          //地址
	Plan            string  `json:"plan"`             //个人计划
	Experience      string  `json:"experience"`       //项目实践经历
	Honor           string  `json:"honor"`            //获得荣誉
	Interest        string  `json:"interest"`         //个人专业研究兴趣方向
	TargetID        int     `json:"target_id"`        //目标导师id
	TargetStatus    int     `json:"target_agree"`     //目标导师状态 		0：初始  1：待处理 2：同意 3：拒绝
	SelectionTable  string  `json:"selection_table"`  //双向选择表url
	AdminStatus     int     `json:"admin_agree"`      //管理员状态			0：初始  1：待处理 2：同意 3：拒绝
	Teacher         Teacher `json:"teacher"`          //导师
	TeacherID       int     `json:"teacher_id"`       //导师id
}
