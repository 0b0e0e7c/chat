package logic

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"chatting/service/userService/internal/svc"
	"chatting/service/userService/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserNameByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserNameByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserNameByIdLogic {
	return &GetUserNameByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserNameByIdLogic) GetUserNameById(in *user.GetUserNameByIdRequest) (*user.GetUserNameByIdResponse, error) {
	userDao := l.svcCtx.UserDao

	names, err := userDao.GetUserNamesFromHash(l.ctx, in.UserId)
	if err != nil {
		logx.Info("User names not found in cache")
	}
	var res = make(map[int64]string)

	var unknown []int64
	for _, uid := range in.UserId {
		if _, ok := names[uid]; ok {
			res[uid] = names[uid]
		} else {
			unknown = append(unknown, uid)
		}
	}

	users, err := userDao.FindUsersByIds(unknown)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	err = userDao.PutUserNameToSet(l.ctx, users)
	if err != nil {
		logx.Errorf("PutUserNameToSet error: %v", err)
	}

	for _, u := range users {
		res[u.UID] = u.Username
	}

	return &user.GetUserNameByIdResponse{
		UserNames: res,
	}, nil
}
