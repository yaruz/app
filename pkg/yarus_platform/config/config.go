package config

import (
	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/db/redis"
	"github.com/minipkg/log"
)

type Configuration struct {
	Log   log.Config
	Infra Infrastructure
}

type Infrastructure struct {
	DataDB        minipkg_gorm.Config
	SearchDB      minipkg_gorm.Config
	Redis         redis.Config
	CacheLifeTime uint
}
