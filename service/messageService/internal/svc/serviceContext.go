package svc

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"chatting/service/messageService/internal/config"
	"chatting/service/messageService/internal/dao"
	"chatting/service/messageService/internal/model"
)

type ServiceContext struct {
	Config config.Config

	MessageDao *dao.MessageDao
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&model.Message{})
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
		Config:     c,
		MessageDao: dao.NewMessageDao(db, dao.WithCache(redisClient)),
	}
}
