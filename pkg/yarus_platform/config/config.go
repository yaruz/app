package config

import (
	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/db/redis"
	"github.com/minipkg/log"
)

type Configuration struct {
	Infrastructure *Infrastructure
	Metadata       *Metadata
}

type Infrastructure struct {
	Log           log.Config
	ReferenceDB   minipkg_gorm.Config
	DataDB        minipkg_gorm.Config
	SearchDB      minipkg_gorm.Config
	Redis         redis.Config
	CacheLifeTime uint
}
