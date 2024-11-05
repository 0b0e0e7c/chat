package logic

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"chatting/service/friendService/internal/model"
	"chatting/service/friendService/internal/svc"
	"chatting/service/friendService/pb/friend"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFriendToGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFriendToGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFriendToGroupLogic {
	return &AddFriendToGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddFriendToGroupLogic) AddFriendToGroup(in *friend.AddFriendToGroupRequest) (*friend.GeneralResponse, error) {
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

	err = friendDao.AddFriendToGroup(in.GetGroupId(), in.GetFriendId())
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			return nil, errors.New("friend already in group")
		}
		return nil, fmt.Errorf("failed to add friend to group: %w", err)
	}

	return &friend.GeneralResponse{
		Success: true,
	}, nil
}
