package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DaoOption func(*FriendDao)

// WithCache sets the cache client in DaoBase
func WithCache(cache *redis.Client) DaoOption {
	return func(dao *FriendDao) {
		dao.cache = cache
	}
}

// applyOptions applies the provided options to the DaoBase
func (d *FriendDao) applyOptions(opts ...DaoOption) {
	for _, opt := range opts {
		opt(d)
	}
}

type FriendDao struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewFriendDao(db *gorm.DB) *FriendDao {
	dao := &FriendDao{
		db: db,
	}

	return dao
}
