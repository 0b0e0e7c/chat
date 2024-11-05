package logic

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"chatting/component/auth"
	"chatting/component/common"
	"chatting/service/userService/internal/svc"
	"chatting/service/userService/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(req *user.LoginRequest) (*user.LoginResponse, error) {
	if req.GetUsername() == "" || req.GetPassword() == "" {
		return nil, errors.New("invalid username or password")
	}

	hashedPassword := common.Hashing(req.Username, req.Password)

	userDao := l.svcCtx.UserDao

	loginUser, err := userDao.FindUserByUsernameAndPassword(req.Username, hashedPassword)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	// 生成JWT Token
	token, err := auth.GenerateToken(int64(loginUser.UID), loginUser.Username, "user-service")
	if err != nil {
		return nil, err
	}

	// 返回登录响应
	resp := &user.LoginResponse{
		UserId:   int64(loginUser.UID),
		Username: loginUser.Username,
		Token:    token,
	}

	return resp, nil
}
