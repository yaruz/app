package config

import (
	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/db/redis"
	"github.com/minipkg/log"
)

type Configuration struct {
	Infrastructure Infrastructure
	ReferenceData  ReferenceData
}

type Infrastructure struct {
	Log           log.Config
	ReferenceDB   minipkg_gorm.Config
	DataDB        minipkg_gorm.Config
	SearchDB      minipkg_gorm.Config
	Redis         redis.Config
	CacheLifeTime uint
}

type ReferenceData struct {
	Languages     Languages
	PropertyUnits PropertyUnits
}

type Languages []Language

type Language struct {
	Code string
	Name string
}

type PropertyUnits []PropertyUnit

type PropertyUnit struct {
	Sysname string
	Texts   map[string]NameAndDescriptionText
}

type NameAndDescriptionText struct {
	Name        string
	Description string
}
