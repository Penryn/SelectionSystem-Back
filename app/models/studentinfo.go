package models

type Student struct {
	ID              int     `json:"id"`               //用户id
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
	Target_id       int     `json:"target_id"`        //目标导师id
	Teacher         Teacher `json:"teacher"`          //导师
}
