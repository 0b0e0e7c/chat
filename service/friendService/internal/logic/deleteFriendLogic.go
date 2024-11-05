package logic

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"chatting/service/friendService/internal/model"
	"chatting/service/friendService/internal/svc"
	"chatting/service/friendService/pb/friend"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFriendLogic {
	return &DeleteFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFriendLogic) DeleteFriend(in *friend.DeleteFriendRequest) (*friend.GeneralResponse, error) {
	friendDao := l.svcCtx.FriendDao

	if in.GetUserId() == 0 || in.GetFriendId() == 0 {
		return nil, ErrUserIdOrFriendIdRequired
	}

	existingFriend, err := friendDao.ValidateFriendship(in.GetUserId(), in.GetFriendId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("operation not allowed, users are not friends")
		}
		return nil, err
	}

	// 已存在好友关系且状态为已添加，则设置为已删除
	if existingFriend.Status == model.Added {
		err := friendDao.UpdateFriendStatus(existingFriend, model.Deleted)
		if err != nil {
			return nil, err
		}
		return &friend.GeneralResponse{
			Success: true,
		}, nil
	}

	return &friend.GeneralResponse{}, nil
}
