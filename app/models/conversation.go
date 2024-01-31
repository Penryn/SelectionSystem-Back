package models

type Conversation struct {
	Id      int64  `json:"id"`       //编号
	UserAID int64  `json:"userA_id"` //用户A的ID
	UserBID int64  `json:"userB_id"` //用户B的ID
	Content string `json:"content"`  //内容
	Time    string `json:"time"`     //时间
}
