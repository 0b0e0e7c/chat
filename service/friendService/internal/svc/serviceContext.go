package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"chatting/service/friendService/internal/config"
	"chatting/service/friendService/internal/dao"
	"chatting/service/friendService/internal/model"
)

type ServiceContext struct {
	Config    config.Config
	FriendDao *dao.FriendDao
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(
		&model.Friend{},
		&model.FriendGroup{},
		&model.FriendGroupMember{},
	)
	if err != nil {
		logx.Errorf("failed to auto migrate: %v", err)
	}

	return &ServiceContext{
		Config:    c,
		FriendDao: dao.NewFriendDao(db),
	}
}
