package logic

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"chatting/service/userService/internal/model"
	"chatting/service/userService/internal/svc"
	"chatting/service/userService/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserProfileLogic {
	return &UpdateUserProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserProfileLogic) UpdateUserProfile(in *user.UpdateUserProfileRequest) (
	*user.UpdateUserProfileResponse, error) {
	if in.Email == "" || in.Nickname == "" || in.Avatar == "" || in.UserId == 0 {
		return nil, fmt.Errorf(
			"invalid input: userId=%d email=%s, nickname=%s, avatar=%s", in.UserId, in.Email, in.Nickname, in.Avatar)
	}

	userDao := l.svcCtx.UserDao

	_, err := userDao.FindUserProfileByUID(in.GetUserId())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// create new profile
		_, err := userDao.CreateUserProfile(
			in.UserId, in.Avatar, in.Nickname, in.Email, model.UserStatus(in.Status))
		if err != nil {
			return nil, err
		}
	} else if err == nil {
		// update profile
		err := userDao.UpdateUserProfile(
			in.UserId, in.Avatar, in.Nickname, in.Email, model.UserStatus(in.Status))
		if err != nil {
			return nil, fmt.Errorf("failed to update user profile: %w", err)
		}
	}

	resp := &user.UpdateUserProfileResponse{
		Success: true,
	}

	return resp, nil
}
