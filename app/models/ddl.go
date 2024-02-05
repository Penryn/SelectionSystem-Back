package models

type DDL struct {
	Id        int  `json:"id"`         //编号
	UserID    string `json:"user_id"`    //用户编号
	DDLType   int    `json:"ddl_type"`   //DDL类型   1:老师DDL  2:管理员DDL
	FirstDDL  string `json:"first_ddl"`  //第一次DDL
	SecondDDL string `json:"second_ddl"` //第二次DDL
}
