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

type ConfirmRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmRequestLogic {
	return &ConfirmRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfirmRequestLogic) ConfirmRequest(in *friend.AddFriendRequest) (*friend.GeneralResponse, error) {
	if in.GetUserId() == 0 || in.GetFriendId() == 0 {
		return nil, ErrUserIdOrFriendIdRequired
	}
	if in.GetUserId() == in.GetFriendId() {
		return nil, ErrCannotAddYourself
	}

	friendDao := l.svcCtx.FriendDao

	// check record exist
	existingFriend, err := friendDao.ValidateFriendship(in.UserId, in.FriendId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果不存在，则创建新的好友关系
			return nil, errors.New("should pending friend request first")
		}
		return nil, errors.New("failed to validate friendship")
	}

	// 已存在好友关系且状态为已添加，则返回错误
	if existingFriend.Status == model.Added {
		return &friend.GeneralResponse{Success: false}, ErrAlreadyFriends
	}

	// 待添加好友的情况，此时更新状态为已添加
	if existingFriend.Status == model.Pending {
		err := friendDao.UpdateFriendStatus(existingFriend, model.Added)
		if err != nil {
			return nil, err
		}
	}

	// 如果已删除，更新状态为已添加
	if existingFriend.Status == model.Deleted {
		err := friendDao.UpdateFriendStatus(existingFriend, model.Added)
		if err != nil {
			return nil, fmt.Errorf("failed to add friend: %w", err)
		}
	}
	return &friend.GeneralResponse{Success: true}, nil
}
