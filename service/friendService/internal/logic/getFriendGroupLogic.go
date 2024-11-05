package logic

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"chatting/service/friendService/internal/svc"
	"chatting/service/friendService/pb/friend"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendGroupLogic {
	return &GetFriendGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendGroupLogic) GetFriendGroup(in *friend.GetFriendGroupRequest) (*friend.GetFriendGroupResponse, error) {
	if in.GetUserId() == 0 {
		return nil, ErrUserIdRequired
	}

	friendDao := l.svcCtx.FriendDao

	// 获取好友组和好友组成员
	friendGroupsWithMembers, err := friendDao.GetFriendGroupsWithMembers(in.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get friend group with members failed: %w", err)
		}
		return nil, err
	}

	// logx.Infof("friendGroupsWithMembers: %+v", friendGroupsWithMembers)

	// 构建响应
	respGroup := make(map[int64]*friend.FriendGroupItem)
	for _, item := range friendGroupsWithMembers {
		if _, ok := respGroup[item.GroupID]; !ok {
			respGroup[item.GroupID] = &friend.FriendGroupItem{
				GroupId:   item.GroupID,
				GroupName: item.GroupName,
				Friends:   make([]*friend.FriendItem, 0),
			}
		}

		// 如果好友信息存在，添加到该组的好友列表中
		if item.FriendID != 0 {
			respGroup[item.GroupID].Friends = append(
				respGroup[item.GroupID].Friends, &friend.FriendItem{
					UserId: item.FriendID,
				})
		}
	}

	groups := make([]*friend.FriendGroupItem, 0, len(respGroup))
	for _, group := range respGroup {
		groups = append(groups, group)
	}

	resp := &friend.GetFriendGroupResponse{
		Groups: groups,
	}

	return resp, nil
}
