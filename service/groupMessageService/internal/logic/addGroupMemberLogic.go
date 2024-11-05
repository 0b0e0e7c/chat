package logic

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"chatting/service/groupMessageService/internal/svc"
	"chatting/service/groupMessageService/pb/groupMessage"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type AddGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddGroupMemberLogic {
	return &AddGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddGroupMemberLogic) AddGroupMember(in *groupMessage.AddGroupMemberRequest) (
	*groupMessage.AddGroupMemberResponse, error) {
	gmDao := l.svcCtx.GroupMessageDao

	err := gmDao.CheckGroupExistById(in.GetGroupId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("group does not exist")
		}
		return nil, err
	}

	// check if the user is the owner
	err = gmDao.CheckOwner(in.GroupId, in.UserId)
	if err != nil {
		return nil, errors.New("user is not the owner of the group")
	}

	err = gmDao.AddGroupMemberWithOwner(in.GetGroupId(), in.GetUserId(), in.GetMemberId())
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062") {
			return nil, errors.New("member already in group")
		}
		return nil, fmt.Errorf("failed to add group member: %v", err)
	}

	return &groupMessage.AddGroupMemberResponse{Success: true}, nil
}
