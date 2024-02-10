package models

type Reason struct {
	ID   int  `json:"id"`
	UserID int `json:"user_id"` //用户id
	ReasonName string `json:"reason_name"` //原因名称
}