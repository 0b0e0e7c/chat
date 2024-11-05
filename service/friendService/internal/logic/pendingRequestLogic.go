package logic

import (
	"context"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"chatting/service/friendService/internal/model"
	"chatting/service/friendService/internal/svc"
	"chatting/service/friendService/pb/friend"
)

type PendingRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPendingRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PendingRequestLogic {
	return &PendingRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PendingRequestLogic) PendingRequest(in *friend.AddFriendRequest) (*friend.GeneralResponse, error) {
	if in.GetUserId() == 0 || in.GetFriendId() == 0 {
		return nil, ErrUserIdOrFriendIdRequired
	}

	if in.GetUserId() == in.GetFriendId() {
		return nil, ErrCannotAddYourself
	}

	friendDao := l.svcCtx.FriendDao

	existingFriend, err := friendDao.ValidateFriendship(in.GetUserId(), in.GetFriendId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果不存在，则创建新的好友关系
			if err := friendDao.PendingFriend(in.UserId, in.FriendId); err != nil {
				return nil, fmt.Errorf("failed when pending friend: %w", err)
			}
		}
		return nil, fmt.Errorf("failed when validate friendship: %w", err)
	}

	// 已存在好友关系且状态为已添加，则返回错误
	if existingFriend.Status == model.Added {
		return &friend.GeneralResponse{Success: false}, ErrAlreadyFriends
	}

	// 如果已删除，更新状态为待添加
	if existingFriend.Status == model.Deleted {
		err := friendDao.UpdateFriendStatus(existingFriend, model.Pending)
		if err != nil {
			return nil, fmt.Errorf("failed when pending friend: %w", err)
		}
	}

	return &friend.GeneralResponse{Success: true}, nil
}
