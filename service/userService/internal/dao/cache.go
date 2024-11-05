package dao

import (
	"context"
	"strconv"
	"time"

	"chatting/service/userService/internal/model"
)

func (dao *UserDao) SetKV(ctx context.Context, k string, v any, expire time.Duration) (err error) {
	return dao.cache.Set(ctx, k, v, expire).Err()
}

func (dao *UserDao) GetKV(ctx context.Context, k string) (v string, err error) {
	return dao.cache.Get(ctx, k).Result()
}

func (dao *UserDao) DelKV(ctx context.Context, k string) {
	dao.cache.Del(ctx, k)
}

// PutUserNameToSet  将用户名存入哈希
func (dao *UserDao) PutUserNameToSet(ctx context.Context, users []*model.User) (err error) {
	if len(users) == 0 {
		return nil
	}

	pipe := dao.cache.Pipeline()
	key := "user_id_name"

	for _, user := range users {
		pipe.HSet(ctx, key, user.UID, user.Username)
	}

	_, err = pipe.Exec(ctx)
	return err
}

// GetUserNamesFromHash 从哈希中获取用户名
func (dao *UserDao) GetUserNamesFromHash(ctx context.Context, uids []int64) (name map[int64]string, err error) {
	key := "user_id_name"
	name = make(map[int64]string)

	for _, uid := range uids {
		userName, err := dao.cache.HGet(ctx, key, strconv.FormatInt(uid, 10)).Result()
		if err != nil {
			return nil, err
		}
		name[uid] = userName
	}

	return name, nil
}
