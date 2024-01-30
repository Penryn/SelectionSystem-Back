package models

type Reason struct {
	Id     int64  `json:"id"`
	Reason string `json:"reason"` //原因
	Advice string `json:"advice"` //理由
}
