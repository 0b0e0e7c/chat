package logic

import (
	"context"
	"fmt"

	"chatting/service/friendService/internal/svc"
	"chatting/service/friendService/pb/friend"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFriendGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateFriendGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFriendGroupLogic {
	return &CreateFriendGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateFriendGroupLogic) CreateFriendGroup(in *friend.CreateFriendGroupRequest) (
	*friend.CreateFriendGroupResponse, error) {
	friendDao := l.svcCtx.FriendDao

	newGroup, err := friendDao.CreateFriendGroup(in.GetUserId(), in.GetGroupName())
	if err != nil {
		return nil, fmt.Errorf("failed to create friend group: %w", err)
	}

	return &friend.CreateFriendGroupResponse{
		GroupId: newGroup.ID,
	}, nil
}
