package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"chatting/service/friendService/internal/model"
	"chatting/service/friendService/internal/svc"
	"chatting/service/friendService/pb/friend"
)

type CheckIsFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckIsFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckIsFriendLogic {
	return &CheckIsFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckIsFriendLogic) CheckIsFriend(in *friend.CheckIsFriendRequest) (*friend.CheckIsFriendResponse, error) {
	if in.GetUserId() == 0 || in.GetFriendId() == 0 {
		return nil, ErrUserIdOrFriendIdRequired
	}
	friendDao := l.svcCtx.FriendDao

	existingFriend, err := friendDao.ValidateFriendship(in.UserId, in.FriendId)
	if err == nil {
		if existingFriend.Status == model.Added {
			return &friend.CheckIsFriendResponse{IsFriend: true}, nil
		}
	}

	return &friend.CheckIsFriendResponse{}, nil
}
