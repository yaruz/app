package redis

import (
	"github.com/minipkg/db/redis"
	goredis "github.com/redis/go-redis/v9"
)

// Repository is an interface of repository
type Repository interface {
	DB() goredis.Cmdable
	Close() error
}

// repository persists albums in database
type repository struct {
	db redis.IDB
}

func NewRepository(db redis.IDB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) DB() goredis.Cmdable {
	return r.db.DB()
}

func (r *repository) Close() error {
	return r.db.Close()
}
