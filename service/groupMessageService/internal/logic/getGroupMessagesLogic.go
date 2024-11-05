package logic

import (
	"context"
	"time"

	"chatting/service/groupMessageService/internal/svc"
	"chatting/service/groupMessageService/pb/groupMessage"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMessagesLogic {
	return &GetGroupMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupMessagesLogic) GetGroupMessages(in *groupMessage.GetGroupMessagesRequest) (
	*groupMessage.GetGroupMessagesResponse, error) {
	if in.Limit == 0 {
		in.Limit = 100
	}

	gmDao := l.svcCtx.GroupMessageDao

	// check if group exists
	err := gmDao.CheckGroupExistById(in.GroupId)
	if err != nil {
		return nil, err
	}

	// check if user is a member of the group
	err = gmDao.CheckExistInGroup(in.GroupId, in.UserId)
	if err != nil {
		return nil, err
	}

	// get messages from Redis
	messages, err := gmDao.PullGroupMsgFromRedis(l.ctx, in.GroupId, in.Offset, in.Limit)
	if err != nil {
		return nil, err
	}

	msgLen := len(messages)
	var latestMsg time.Time
	if msgLen > 0 {
		latestMsg = messages[msgLen-1].Timestamp
	}

	// get messages from MySQL
	if msgLen == 0 || msgLen < int(in.Limit) {
		messages, err = gmDao.PullGroupMsgFromMySQL(in.GroupId, in.Offset, in.Limit)
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
			if err := gmDao.PushGroupMsgToRedis(l.ctx, newMessages); err != nil {
				return nil, err
			}
		}
	}

	msgList := make([]*groupMessage.GroupMessage, 0, len(messages))
	for _, msg := range messages {
		msgList = append(
			msgList, &groupMessage.GroupMessage{
				MsgId:     msg.MsgID,
				SenderId:  msg.SenderId,
				GroupId:   msg.GroupId,
				Content:   msg.Content,
				Timestamp: msg.Timestamp.Unix(),
			})
	}

	return &groupMessage.GetGroupMessagesResponse{
		Messages: msgList,
	}, nil
}
