package model

import (
	"time"

	"gorm.io/gorm"
)

type ChatGroup struct {
	GID     int64  `gorm:"primaryKey;autoIncrement;type:bigint;column:gid"`
	Name    string `gorm:"uniqueIndex;size:64"`        // 群组名称
	Avatar  string `gorm:"size:256"`                   // 群组头像 URL
	OwnerID int64  `gorm:"index;type:bigint;not null"` // 群主 ID

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// 群组成员中间表
type ChatGroupMember struct {
	GroupID int64 `gorm:"primaryKey;type:bigint;not null"` // 将 GroupID 设为主键的一部分
	UserID  int64 `gorm:"primaryKey;type:bigint;not null"` // 将 UserID 设为主键的一部分

	CreatedAt time.Time
	UpdatedAt time.Time
}
