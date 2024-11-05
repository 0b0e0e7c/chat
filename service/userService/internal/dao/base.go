package dao

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserDao struct {
	db    *gorm.DB
	cache *redis.Client
}

type DaoOption func(*UserDao)

// WithCache sets the cache client in DaoBase
func WithCache(cache *redis.Client) DaoOption {
	return func(dao *UserDao) {
		dao.cache = cache
	}
}

// applyOptions applies the provided options to the DaoBase
func (d *UserDao) applyOptions(opts ...DaoOption) {
	for _, opt := range opts {
		opt(d)
	}
}

func NewUserDao(db *gorm.DB, opts ...DaoOption) *UserDao {
	dao := &UserDao{
		db: db,
	}

	dao.applyOptions(opts...)

	return dao
}
