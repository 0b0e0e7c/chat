package logic

import (
	"context"
	"fmt"
	"strings"

	"chatting/component/common"
	"chatting/service/userService/internal/svc"
	"chatting/service/userService/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	if in.Username == "" || in.Password == "" {
		return nil, fmt.Errorf("username or password is empty")
	}

	userDao := l.svcCtx.UserDao
	newUser, err := userDao.CreateUserByUsernameAndPassword(in.Username, common.Hashing(in.Username, in.Password))
	if err != nil {
		if strings.HasPrefix(err.Error(), "Error 1062") {
			return nil, fmt.Errorf("username already exists")
		}
		return nil, fmt.Errorf("create user fail: %v", err)
	}

	resp := &user.RegisterResponse{
		UserId: newUser.UID,
	}

	return resp, nil
}
