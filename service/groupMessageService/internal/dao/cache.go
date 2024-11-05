package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"chatting/service/groupMessageService/internal/model"
)

// PushGroupMsgToRedis push messages to Redis
func (dao *GroupMessageDao) PushGroupMsgToRedis(ctx context.Context, messages []*model.GroupMessage) error {
	if len(messages) == 0 {
		return nil
	}

	key := fmt.Sprintf("group:%d", messages[0].GroupId)

	pipe := dao.cache.Pipeline()

	for _, msg := range messages {
		msgEntry, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		logx.Infof("put group message to Redis: %+v", msg)

		pipe.ZAdd(
			ctx, key, redis.Z{
				Score:  float64(msg.MsgID),
				Member: msgEntry,
			})
	}

	pipe.ZRemRangeByRank(ctx, key, 0, -11)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to push new messages to redis: %v", err)
	}
	return nil
}

// PullGroupMsgFromRedis pull messages from Redis
func (dao *GroupMessageDao) PullGroupMsgFromRedis(
	ctx context.Context, groupID, offset, limit int64) ([]*model.GroupMessage, error) {
	logx.Infof("pull group %d 's messages from Redis", groupID)

	key := fmt.Sprintf("group:%d", groupID)

	msgs, err := dao.cache.ZRange(ctx, key, offset, limit).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to pull messages to redis: %v", err)
	}

	messages := make([]*model.GroupMessage, len(msgs))
	for _, msgStr := range msgs {
		var msg *model.GroupMessage
		if err := json.Unmarshal([]byte(msgStr), &msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	logx.Infof("get %d messages from Redis", len(messages))
	return messages, nil
}
