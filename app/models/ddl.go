package models

type DDL struct {
	Id        int64  `json:"id"`		 //编号
	User      string `json:"user"`      //用户名
	FirstDDL  string `json:"firstDDL"`  //第一次DDL
	SecondDDL string `json:"secondDDL"` //第二次DDL
}
