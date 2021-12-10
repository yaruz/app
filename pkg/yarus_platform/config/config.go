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
	DataSharding  Sharding
	SearchDB      minipkg_gorm.Config
	Redis         redis.Config
	CacheLifeTime uint
}

type Sharding struct {
	IsAutoMigrate bool
	Default       DBCluster
	BySysnames    map[string]DBCluster
}

type DBCluster struct {
	Capacity    uint
	EntityTypes []string
	Items       []minipkg_gorm.Config
}

func (s *Sharding) GetEntityTypesByClusterSysnames() map[string][]string {
	types := make(map[string][]string, len(s.BySysnames))

	for sysname, cluster := range s.BySysnames {
		types[sysname] = cluster.EntityTypes
	}
	return types
}
