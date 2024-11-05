package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"chatting/service/groupMessageService/internal/model"
	"chatting/service/groupMessageService/internal/svc"
	"chatting/service/groupMessageService/pb/groupMessage"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendGroupMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendGroupMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendGroupMessageLogic {
	return &SendGroupMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendGroupMessageLogic) SendGroupMessage(in *groupMessage.SendGroupMessageRequest) (
	*groupMessage.SendGroupMessageResponse, error) {
	gmDao := l.svcCtx.GroupMessageDao

	err := gmDao.CheckGroupExistById(in.GetGroupId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("group does not exist")
		}
		return nil, err
	}

	err = gmDao.CheckExistInGroup(in.GetGroupId(), in.GetSenderId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cannot send message, user not in group")
		}
		return nil, err
	}

	newGroupMsg := &model.GroupMessage{
		GroupId:   in.GetGroupId(),
		SenderId:  in.GetSenderId(),
		Content:   in.GetContent(),
		Timestamp: time.Now(),
	}

	if err = gmDao.CreateGroupMessage(newGroupMsg); err != nil {
		return nil, fmt.Errorf("failed to create group message: %v", err)
	}

	if err = gmDao.PushGroupMsgToRedis(l.ctx, []*model.GroupMessage{newGroupMsg}); err != nil {
		return nil, fmt.Errorf("failed to push group message to redis: %v", err)
	}

	return &groupMessage.SendGroupMessageResponse{
		Success:   true,
		MsgId:     newGroupMsg.MsgID,
		Timestamp: newGroupMsg.Timestamp.Unix(),
	}, nil
}
