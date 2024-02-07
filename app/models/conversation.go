package models

type Conversation struct {
	ID      int  `json:"id"`       //编号
	UserAID int  `json:"userA_id"` //用户A的ID
	UserBID int  `json:"userB_id"` //用户B的ID
	Content string `json:"content"`  //内容
	Time    string `json:"time"`     //时间
}
