package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DaoOption func(*GroupMessageDao)

// WithCache sets the cache client in DaoBase
func WithCache(cache *redis.Client) DaoOption {
	return func(dao *GroupMessageDao) {
		dao.cache = cache
	}
}

// applyOptions applies the provided options to the DaoBase
func (d *GroupMessageDao) applyOptions(opts ...DaoOption) {
	for _, opt := range opts {
		opt(d)
	}
}

type GroupMessageDao struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewMessageGroupDao(db *gorm.DB, opts ...DaoOption) *GroupMessageDao {
	dao := &GroupMessageDao{
		db: db,
	}

	dao.applyOptions(opts...)
	return dao
}
