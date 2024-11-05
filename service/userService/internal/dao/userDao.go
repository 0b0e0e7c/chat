package dao

import (
	"chatting/service/userService/internal/model"
)

func (dao *UserDao) CreateUserByUsernameAndPassword(username, password string) (newUser *model.User, err error) {
	newUser = &model.User{
		Username: username,
		Password: password,
	}

	return newUser, dao.db.Create(&newUser).Error
}

func (dao *UserDao) FindUserByUsernameAndPassword(username, password string) (*model.User, error) {
	var loginUser *model.User

	return loginUser, dao.db.Where("username = ? AND password = ?", username, password).First(&loginUser).Error
}

func (dao *UserDao) FindUserById(id int64) (*model.User, error) {
	var user *model.User

	return user, dao.db.Where("uid = ?", id).First(&user).Error
}

func (dao *UserDao) FindUsersByIds(ids []int64) ([]*model.User, error) {
	var users []*model.User

	return users, dao.db.Where("uid IN ?", ids).Find(&users).Error
}

func (dao *UserDao) UpdateUserProfile(uid int64, avatar, nickname, email string, status model.UserStatus) error {
	return dao.db.Model(&model.UserProfile{}).Where("uid = ?", uid).Updates(
		model.UserProfile{
			Avatar:   avatar,
			Nickname: nickname,
			Email:    email,
			Status:   status,
		}).Error
}

func (dao *UserDao) CreateUserProfile(
	uid int64, avatar, nickname, email string, status model.UserStatus) (newProfile *model.UserProfile, err error) {
	newProfile = &model.UserProfile{
		UID:      uid,
		Avatar:   avatar,
		Nickname: nickname,
		Email:    email,
		Status:   status,
	}
	return newProfile, dao.db.Create(&newProfile).Error
}

func (dao *UserDao) FindUserProfileByUID(uid int64) (profile *model.UserProfile, err error) {
	return profile, dao.db.Where("uid = ?", uid).First(&profile).Error
}
