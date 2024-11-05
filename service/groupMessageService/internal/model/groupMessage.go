package model

import "time"

// group message
type GroupMessage struct {
	MsgID     int64     `gorm:"primaryKey;autoIncrement" json:"msg_id"`
	SenderId  int64     `gorm:"index;not null" json:"sender_id"`
	GroupId   int64     `gorm:"index;not null" json:"group_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Timestamp time.Time `gorm:"not null" json:"timestamp"`
}
