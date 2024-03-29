package models

import "time"

type DDL struct {
	ID        int    `json:"id"`         //编号
	UserID    int `json:"user_id"`    //用户编号
	DDLType   int    `json:"ddl_type"`   //DDL类型   1:老师DDL  2:管理员DDL
	FirstDDL  time.Time `json:"first_ddl"`  //第一次DDL
	SecondDDL time.Time `json:"second_ddl"` //第二次DDL
}
