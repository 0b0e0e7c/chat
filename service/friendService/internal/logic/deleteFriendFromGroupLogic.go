package logic

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"chatting/service/friendService/internal/model"
	"chatting/service/friendService/internal/svc"
	"chatting/service/friendService/pb/friend"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFriendFromGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFriendFromGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFriendFromGroupLogic {
	return &DeleteFriendFromGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFriendFromGroupLogic) DeleteFriendFromGroup(in *friend.AddFriendToGroupRequest) (
	*friend.GeneralResponse, error) {
	if in.GetUserId() == 0 || in.GetFriendId() == 0 || in.GetGroupId() == 0 {
		return nil, errors.New("user id, friend id and group id are required")
	}

	friendDao := l.svcCtx.FriendDao

	existingFriend, err := friendDao.ValidateFriendship(in.GetUserId(), in.GetFriendId())
	if errors.Is(err, gorm.ErrRecordNotFound) || (err == nil && existingFriend.Status != model.Added) {
		return nil, errors.New("operation not allowed, users are not friends")
	}
	if err != nil {
		return nil, err
	}

	err = friendDao.CheckFriendGroupExist(in.GetUserId(), in.GetGroupId())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("group not found")
	}

	err = friendDao.CheckFriendExistInGroup(in.GetGroupId(), in.GetFriendId())
	if err != nil {
		return nil, errors.New("friend not in group")
	}

	err = friendDao.DeleteFriendFromGroup(in.GetGroupId(), in.GetFriendId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete friend from group: %w", err)
	}

	return &friend.GeneralResponse{
		Success: true,
	}, nil
}
