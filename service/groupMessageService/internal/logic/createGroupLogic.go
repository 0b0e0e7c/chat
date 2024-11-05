package logic

import (
	"context"
	"errors"
	"fmt"

	"chatting/service/groupMessageService/internal/svc"
	"chatting/service/groupMessageService/pb/groupMessage"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateGroupLogic) CreateGroup(in *groupMessage.CreateGroupRequest) (*groupMessage.CreateGroupResponse, error) {
	if in.GroupName == "" {
		return nil, errors.New("group name can not be empty")
	}

	gmDao := l.svcCtx.GroupMessageDao

	err := gmDao.CheckGroupExistByName(in.GetGroupName())
	if err == nil {
		return nil, errors.New("group name already exists")
	}

	newGroup, err := gmDao.CreateMessageGroup(in.GetGroupName(), in.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %v", err)
	}

	resp := &groupMessage.CreateGroupResponse{
		GroupId: newGroup.GID,
		Success: true,
	}

	return resp, nil
}
