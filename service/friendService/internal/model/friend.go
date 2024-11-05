package model

import (
	"time"

	"gorm.io/gorm"
)

// Friend 表结构
type Friend struct {
	gorm.Model
	UserID      int64        `gorm:"index;not null" json:"user_id"`
	FriendID    int64        `gorm:"index;not null" json:"friend_id"`
	Status      FriendStatus `gorm:"not null" json:"status"`
	InitiatorID int64        `gorm:"not null" json:"initiator_id"`
}

// FriendGroup 用户创建的好友分组
type FriendGroup struct {
	ID   int64  `gorm:"primaryKey;autoIncrement;type:bigint"`
	UID  int64  `gorm:"index;type:bigint;not null"` // 分组所属用户
	Name string `gorm:"size:64;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// GroupMembers represents the members of the friend group.
	// It uses the GroupID field in the FriendGroupMember struct as the foreign key
	// and references the ID field in the FriendGroup struct.
	GroupMembers []FriendGroupMember `gorm:"foreignKey:GroupID;references:ID"`
}

// FriendGroupMember 好友分组成员
type FriendGroupMember struct {
	GroupID  int64 `gorm:"primaryKey;type:bigint;not null"` // 分组 ID
	FriendID int64 `gorm:"primaryKey;type:bigint;not null"` // 好友 ID

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type FriendStatus uint8

const (
	Unadded FriendStatus = iota // 0: 未添加
	Added                       // 1: 已添加
	Pending                     // 2: 等待接受
	Deleted                     // 3: 已删除
)
