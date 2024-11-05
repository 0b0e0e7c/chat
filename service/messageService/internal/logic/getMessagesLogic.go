package logic

import (
	"context"
	"fmt"
	"time"

	"chatting/service/messageService/internal/svc"
	"chatting/service/messageService/pb/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessagesLogic {
	return &GetMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessagesLogic) GetMessages(in *message.GetMessagesRequest) (*message.GetMessagesResponse, error) {
	if in.GetUserId() == 0 || in.GetPeerId() == 0 {
		return nil, fmt.Errorf("userId and peerId are required")
	}

	if in.Limit == 0 {
		in.Limit = 100
	}

	messageDao := l.svcCtx.MessageDao

	// 从 Redis 获取消息
	messages, err := messageDao.PullMsgFromRedis(l.ctx, in.UserId, in.PeerId, in.Offset, in.Offset+in.Limit-1)
	if err != nil {
		return nil, err
	}

	msgLen := len(messages)
	var latestMsg time.Time
	if msgLen > 0 {
		latestMsg = messages[msgLen-1].Timestamp
	}

	// get messages from MySQL
	if len(messages) == 0 || len(messages) < int(in.Limit) {
		messages, err = messageDao.PullMsgFromMySQL(in.UserId, in.PeerId, in.Offset, in.Limit)
		if err != nil {
			return nil, err
		}

		// cache messages with newer timestamp than "latestMsg" to Redis
		if len(messages) != 0 {
			cnt := 0
			for _, msg := range messages {
				if msg.Timestamp.After(latestMsg) {
					cnt++
				}
			}
			newMessages := messages[:cnt]
			// cache messages to Redis
			if err := messageDao.PushMsgToRedis(l.ctx, newMessages); err != nil {
				return nil, err
			}
		}
	}

	msgList := make([]*message.Message, 0, len(messages))
	for _, msg := range messages {
		msgList = append(
			msgList, &message.Message{
				MsgId:      msg.MsgID,
				SenderId:   msg.SenderId,
				ReceiverId: msg.ReceiverId,
				Content:    msg.Content,
				Timestamp:  msg.Timestamp.Unix(),
			})
	}

	return &message.GetMessagesResponse{
		Messages: msgList,
	}, nil
}
