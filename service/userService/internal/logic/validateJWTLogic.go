package logic

import (
	"context"
	"errors"

	"chatting/component/auth"
	"chatting/service/userService/internal/svc"
	"chatting/service/userService/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateJWTLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewValidateJWTLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateJWTLogic {
	return &ValidateJWTLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ValidateJWTLogic) ValidateJWT(in *user.ValidateRequest) (*user.ValidateResponse, error) {
	if in.GetToken() == "" {
		return nil, errors.New("token is required")
	}

	userDao := l.svcCtx.UserDao

	// 检查token是否在黑名单中
	v, err := userDao.GetKV(l.ctx, in.Token)
	if err == nil && v == "blacklisted" {
		l.Logger.Info("token is blacklisted, user is logged out")
		return &user.ValidateResponse{
			Valid:  false,
			UserId: 0,
		}, nil
	}

	valid, uid, err := auth.ValidateToken(in.Token)
	if errors.Is(err, auth.ErrTokenExpired) {
		l.Logger.Info("token expired")
	}

	return &user.ValidateResponse{
		Valid:  valid,
		UserId: uid,
	}, err
}
