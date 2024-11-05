package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	UID      int64  `gorm:"primaryKey;autoIncrement;type:bigint;index"`
	Username string `gorm:"uniqueIndex;size:64"`
	Password string `gorm:"size:128"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Profile UserProfile `gorm:"foreignKey:UID;references:UID"` // 关联 UserProfile 表
}

type UserProfile struct {
	gorm.Model
	UID      int64      `gorm:"uniqueIndex;type:bigint;not null"`
	Avatar   string     `gorm:"size:256"`  // 用户头像 URL
	Nickname string     `gorm:"size:64"`   // 用户昵称
	Email    string     `gorm:"size:128"`  // 用户邮箱
	Status   UserStatus `gorm:"default:1"` // 用户状态
}

type UserStatus int32

const (
	Online  UserStatus = iota // 0: 在线
	Offline                   // 1:  离线
)
