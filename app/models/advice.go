package models

import "time"

type Advice struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Content    string    `json:"content"`
	Anonymity  bool      `json:"anonymity"`
	CreateTime time.Time `json:"create_time"`
}
