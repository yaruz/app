package infrastructure

import (
	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/db/redis"
	"github.com/minipkg/log"
	yaruz_config "github.com/yaruz/app/pkg/yarus_platform/config"
)

type Configuration struct {
	Log           log.Config
	CacheLifeTime uint
	DB            *DBConf
}

type DBConf struct {
	Identity     minipkg_gorm.Config
	Reference    minipkg_gorm.Config
	DataSharding yaruz_config.Sharding
	Search       minipkg_gorm.Config
	Redis        redis.Config
}

func (c *Configuration) YaruzConfig(yaruzMetadata *yaruz_config.Metadata) *yaruz_config.Configuration {
	return &yaruz_config.Configuration{
		Infrastructure: &yaruz_config.Infrastructure{
			Log:           c.Log,
			ReferenceDB:   c.DB.Reference,
			DataSharding:  c.DB.DataSharding,
			SearchDB:      c.DB.Search,
			Redis:         c.DB.Redis,
			CacheLifeTime: c.CacheLifeTime,
		},
		Metadata: yaruzMetadata,
	}
}
