package logic

import (
	"context"
	"errors"
	"time"

	"chatting/component/auth"
	"chatting/service/userService/internal/svc"
	"chatting/service/userService/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutLogic) Logout(in *user.LogoutRequest) (*user.LogoutResponse, error) {
	if in.Token == "" {
		return nil, errors.New("token is required")
	}

	userDao := l.svcCtx.UserDao

	// blacklistTLL 至少是 token 的剩余有效时间
	var blacklistTTL time.Duration
	tokenExpTime, err := auth.GetTokenExpiry(in.Token)
	if err != nil {
		l.Logger.Error("failed to get token expiry time:", err)
		blacklistTTL = auth.TokenExpireDuration
	} else {
		blacklistTTL = time.Until(tokenExpTime)
		if blacklistTTL < 0 {
			blacklistTTL = 0
		}
	}

	err = userDao.SetKV(l.ctx, in.Token, "blacklisted", blacklistTTL)
	if err != nil {
		l.Logger.Error("failed to add token to blacklist:", err)
		return &user.LogoutResponse{
			Success: false,
		}, err
	}

	return &user.LogoutResponse{
		Success: true,
	}, nil
}
