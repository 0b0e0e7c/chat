package dao

import (
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"chatting/service/groupMessageService/internal/model"
)

// CreateMessageGroup create a new group
func (dao *GroupMessageDao) CreateMessageGroup(name string, ownerID int64) (*model.ChatGroup, error) {
	var group *model.ChatGroup

	err := dao.db.Transaction(
		func(tx *gorm.DB) error {
			group = &model.ChatGroup{
				Name:    name,
				OwnerID: ownerID,
			}

			if err := tx.Create(&group).Error; err != nil {
				return err
			}

			// 将所有者添加为群组成员
			if err := dao.AddGroupMemberWithTx(tx, group.GID, ownerID); err != nil {
				return err
			}

			return nil
		})
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (dao *GroupMessageDao) AddGroupMemberWithTx(tx *gorm.DB, groupID, userID int64) error {
	member := &model.ChatGroupMember{
		GroupID: groupID,
		UserID:  userID,
	}

	err := tx.Where("gid = ?", groupID).First(&model.ChatGroup{}).Error
	if err != nil {
		return errors.New("group not found")
	}

	err = tx.Where("group_id = ? AND user_id = ?", groupID, userID).First(&model.ChatGroupMember{}).Error
	if err == nil {
		return errors.New("member already exists in the group")
	}

	return tx.Create(&member).Error
}

// AddGroupMemberWithOwner add a member to the group with the owner
func (dao *GroupMessageDao) AddGroupMemberWithOwner(groupID, userID, memberID int64) error {
	member := &model.ChatGroupMember{
		GroupID: groupID,
		UserID:  memberID,
	}

	return dao.db.Create(&member).Error
}

// CheckOwner check if the user is the owner of the group
func (dao *GroupMessageDao) CheckOwner(groupID, userID int64) error {
	group := &model.ChatGroup{}

	return dao.db.Where("gid = ? AND owner_id = ?", groupID, userID).First(&group).Error
}

// CheckGroupExistById check if a group exists
func (dao *GroupMessageDao) CheckGroupExistById(groupID int64) error {
	group := &model.ChatGroup{}

	return dao.db.Where("gid = ?", groupID).First(&group).Error
}

// CheckGroupExistByName check if a group exists by name
func (dao *GroupMessageDao) CheckGroupExistByName(groupName string) error {
	group := &model.ChatGroup{}

	return dao.db.Where("name = ?", groupName).First(&group).Error
}

// CheckExistInGroup  check if a user exists in a group
func (dao *GroupMessageDao) CheckExistInGroup(groupID, userID int64) error {
	member := &model.ChatGroupMember{}

	return dao.db.Where("group_id = ? AND user_id = ?", groupID, userID).First(&member).Error
}

// CreateGroupMessage create a new message in a group
func (dao *GroupMessageDao) CreateGroupMessage(newGroupMsg *model.GroupMessage) error {
	return dao.db.Create(newGroupMsg).Error
}

// PullGroupMsgFromMySQL pull messages from MySQL
func (dao *GroupMessageDao) PullGroupMsgFromMySQL(groupID, offset, limit int64) ([]*model.GroupMessage, error) {
	logx.Infof("pull messages from MySQL")

	var messages []*model.GroupMessage
	err := dao.db.Where("group_id = ?", groupID).
		Order("timestamp DESC").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&messages).Error

	return messages, err
}
