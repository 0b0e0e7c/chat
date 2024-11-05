package logic

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"chatting/service/userService/internal/svc"
	"chatting/service/userService/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserProfileLogic {
	return &GetUserProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserProfileLogic) GetUserProfile(in *user.GetUserProfileRequest) (*user.GetUserProfileResponse, error) {
	if in.UserId == 0 {
		return nil, errors.New("user id is required")
	}

	userDao := l.svcCtx.UserDao
	userProfile, err := userDao.FindUserProfileByUID(in.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	resp := &user.GetUserProfileResponse{
		Nickname: userProfile.Nickname,
		Email:    userProfile.Email,
		Avatar:   userProfile.Avatar,
		Status:   int32(userProfile.Status),
	}

	return resp, nil
}
