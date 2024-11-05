package logic

import (
	"context"
	"fmt"
	"time"

	"chatting/service/messageService/internal/model"
	"chatting/service/messageService/internal/svc"
	"chatting/service/messageService/pb/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMessageLogic) SendMessage(in *message.SendMessageRequest) (*message.SendMessageResponse, error) {
	if in.GetSenderId() == 0 || in.GetReceiverId() == 0 || in.GetContent() == "" {
		return nil, fmt.Errorf("senderId, receiverId and content are required")
	}

	messageDao := l.svcCtx.MessageDao

	// 创建消息
	newMsg := &model.Message{
		SenderId:   in.SenderId,
		ReceiverId: in.ReceiverId,
		Content:    in.Content,
		Timestamp:  time.Now(),
	}

	if err := messageDao.CreateMessage(newMsg); err != nil {
		return nil, fmt.Errorf("failed to create message: %v", err)
	}

	// 推送消息到 Redis
	if err := messageDao.PushMsgToRedis(l.ctx, []*model.Message{newMsg}); err != nil {
		return nil, fmt.Errorf("failed to push message to Redis: %v", err)
	}

	return &message.SendMessageResponse{
		Success:   true,
		MsgId:     newMsg.MsgID,
		Timestamp: newMsg.Timestamp.Unix(),
	}, nil
}
