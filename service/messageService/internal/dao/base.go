package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DaoOption func(*MessageDao)

// WithCache sets the cache client in DaoBase
func WithCache(cache *redis.Client) DaoOption {
	return func(dao *MessageDao) {
		dao.cache = cache
	}
}

// applyOptions applies the provided options to the DaoBase
func (d *MessageDao) applyOptions(opts ...DaoOption) {
	for _, opt := range opts {
		opt(d)
	}
}

type MessageDao struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewMessageDao(db *gorm.DB, opts ...DaoOption) *MessageDao {
	dao := &MessageDao{
		db: db,
	}

	dao.applyOptions(opts...)
	return dao
}
