package model

import (
	"time"
)

// Message 表结构
type Message struct {
	MsgID      int64     `gorm:"primaryKey;autoIncrement" json:"msg_id"`
	SenderId   int64     `gorm:"index;not null" json:"sender_id"`
	ReceiverId int64     `gorm:"index;not null" json:"receiver_id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	Timestamp  time.Time `gorm:"not null" json:"timestamp"`
}
