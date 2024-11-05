package svc

import (
	"chatting/service/groupMessageService/internal/config"
	"chatting/service/groupMessageService/internal/dao"
	"chatting/service/groupMessageService/internal/model"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config          config.Config
	GroupMessageDao *dao.GroupMessageDao
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&model.ChatGroup{}, &model.ChatGroupMember{}, &model.GroupMessage{})
	if err != nil {
		logx.Errorf("failed to auto migrate: %v", err)
		return nil
	}

	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     c.CacheRedis.Host,
			Password: c.CacheRedis.Pass,
		},
	)

	return &ServiceContext{
		Config:          c,
		GroupMessageDao: dao.NewMessageGroupDao(db, dao.WithCache(redisClient)),
	}
}
