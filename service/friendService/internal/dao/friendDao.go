package dao

import (
	"chatting/component/common"
	"chatting/service/friendService/internal/model"
)

// ValidateFriendship checks if the friend relationship already exists
func (dao *FriendDao) ValidateFriendship(userID, friendID int64) (model.Friend, error) {
	var existingFriend model.Friend
	lowerID, higherID := common.LowHigh(userID, friendID)

	return existingFriend, dao.db.Where(
		"user_id = ? AND friend_id = ?", lowerID, higherID).First(&existingFriend).Error
}

// UpdateFriendStatus updates the status of an existing friend relationship
func (dao *FriendDao) UpdateFriendStatus(friend model.Friend, status model.FriendStatus) error {
	return dao.db.Model(&friend).Update("status", status).Error
}

// CreateFriend creates a new friend relationship
func (dao *FriendDao) CreateFriend(userID, friendID int64, status model.FriendStatus) error {
	lowerID, higherID := common.LowHigh(userID, friendID)

	newFriend := model.Friend{
		UserID:      lowerID,
		FriendID:    higherID,
		Status:      status,
		InitiatorID: userID,
	}
	return dao.db.Create(&newFriend).Error
}

// PendingFriend creates a new friend relationship with status pending
func (dao *FriendDao) PendingFriend(userID, friendID int64) error {
	return dao.CreateFriend(userID, friendID, model.Pending)
}

// GetFriendsWithStatus returns a list of friends with the given status
func (dao *FriendDao) GetFriendsWithStatus(userID int64, status model.FriendStatus) ([]model.Friend, error) {
	var friends []model.Friend

	return friends, dao.db.Where("(user_id = ? OR friend_id = ?) AND status = ?", userID, userID, status).
		Find(&friends).Error
}

// CreateFriendGroup creates a new friend group
func (dao *FriendDao) CreateFriendGroup(userID int64, name string) (model.FriendGroup, error) {
	newGroup := model.FriendGroup{
		UID:  userID,
		Name: name,
	}

	return newGroup, dao.db.Create(&newGroup).Error
}

func (dao *FriendDao) CheckFriendGroupExist(userID, groupID int64) error {
	return dao.db.Where("uid = ? AND id = ?", userID, groupID).First(&model.FriendGroup{}).Error
}

func (dao *FriendDao) DeleteFriendGroup(groupID int64) error {
	return dao.db.Where("id = ?", groupID).Delete(&model.FriendGroup{}).Error
}

// AddFriendToGroup adds a friend to a friend group
func (dao *FriendDao) AddFriendToGroup(groupID, friendID int64) error {
	newMember := model.FriendGroupMember{
		GroupID:  groupID,
		FriendID: friendID,
	}
	return dao.db.Create(&newMember).Error
}

// DeleteFriendFromGroup deletes a friend from a friend group
func (dao *FriendDao) DeleteFriendFromGroup(groupID, friendID int64) error {
	// Unscoped to delete the record
	// otherwise there still be a record in the database
	// which we don't need in the case of friend group management
	return dao.db.Unscoped().Where(
		"group_id = ? AND friend_id = ?", groupID, friendID).Delete(&model.FriendGroupMember{}).Error
}

// CheckFriendExistInGroup checks if a friend is in a group
func (dao *FriendDao) CheckFriendExistInGroup(groupID, friendID int64) error {
	return dao.db.Where(
		"group_id = ? AND friend_id = ?", groupID, friendID).First(&model.FriendGroupMember{}).Error
}

// GetFriendGroups returns a list of friend groups
func (dao *FriendDao) GetFriendGroups(userID int64) ([]model.FriendGroup, error) {
	var groups []model.FriendGroup

	return groups, dao.db.Where("uid = ?", userID).Find(&groups).Error
}

// FriendGroupWithMembers represents a friend group with its members
type FriendGroupWithMembers struct {
	GroupID   int64  `gorm:"column:group_id"`   // 好友组ID
	GroupName string `gorm:"column:group_name"` // 好友组名称
	FriendID  int64  `gorm:"column:friend_id"`  // 好友ID
}

// GetFriendGroupsWithMembers returns a list of friend groups with their members
func (dao *FriendDao) GetFriendGroupsWithMembers(userID int64) ([]FriendGroupWithMembers, error) {
	var result []FriendGroupWithMembers

	if err := dao.db.Table("friend_groups").
		Select("friend_groups.id as group_id, friend_groups.name as group_name, friend_group_members.friend_id").
		Joins("LEFT JOIN friend_group_members ON friend_groups.id = friend_group_members.group_id").
		Where("friend_groups.uid = ?", userID).
		Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
