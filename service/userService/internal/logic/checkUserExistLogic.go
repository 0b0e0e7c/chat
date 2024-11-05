package logic

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"chatting/service/userService/internal/svc"
	"chatting/service/userService/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckUserExistLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckUserExistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckUserExistLogic {
	return &CheckUserExistLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckUserExistLogic) CheckUserExist(in *user.CheckUserExistRequest) (*user.CheckUserExistResponse, error) {
	userDao := l.svcCtx.UserDao

	users := in.UserId

	if len(users) == 0 {
		return &user.CheckUserExistResponse{
			Exist: false,
		}, nil
	}

	_, err := userDao.FindUsersByIds(users)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user(s) not found")
		}
		return nil, err
	}

	return &user.CheckUserExistResponse{
		Exist: true,
	}, nil
}
