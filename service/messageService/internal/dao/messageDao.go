package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"chatting/component/common"
	"chatting/service/messageService/internal/model"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

func (dao *MessageDao) PullMsgFromMySQL(senderId, receiverId, offset, limit int64) ([]*model.Message, error) {
	logx.Infof("pull messages from MySQL")

	var messages []*model.Message
	err := dao.db.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		senderId, receiverId, receiverId, senderId).
		Order("timestamp DESC").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&messages).Error

	return messages, err
}

func (dao *MessageDao) PullMsgFromRedis(ctx context.Context, UserID, PeerId, start, end int64) (
	[]*model.Message, error) {

	lowerID, higherID := common.LowHigh(UserID, PeerId)
	key := fmt.Sprintf("chat:%d-%d", lowerID, higherID)
	logx.Infof("pull messages from Redis, chat:%d-%d", lowerID, higherID)

	msgs, err := dao.cache.ZRevRange(ctx, key, start, end).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to pull messages from Redis: %v", err)
	}

	messages := make([]*model.Message, len(msgs))
	for _, msgStr := range msgs {
		var msg *model.Message
		if err := json.Unmarshal([]byte(msgStr), &msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	logx.Infof("messages from Redis: %+v len: %d", messages, len(messages))
	return messages, nil
}

func (dao *MessageDao) CreateMessage(newMsg *model.Message) error {
	return dao.db.Create(newMsg).Error
}

// PushMsgToRedis pushes a list of messages or a single message to Redis.
func (dao *MessageDao) PushMsgToRedis(ctx context.Context, messages []*model.Message) error {
	if len(messages) == 0 {
		return nil
	}

	lowerID, higherID := common.LowHigh(messages[0].SenderId, messages[0].ReceiverId)
	key := fmt.Sprintf("chat:%d-%d", lowerID, higherID)

	pipe := dao.cache.Pipeline()

	for _, msg := range messages {
		msgEntry, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		logx.Infof("put message to Redis: %+v", msg)

		pipe.ZAdd(
			ctx, key, redis.Z{
				Score:  float64(msg.MsgID),
				Member: msgEntry,
			})
	}

	// 保留最新的10条消息
	pipe.ZRemRangeByRank(ctx, key, 0, -11).Err()

	_, err := pipe.Exec(ctx)
	return err
}
