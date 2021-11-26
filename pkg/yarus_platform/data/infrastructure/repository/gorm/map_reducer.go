package gorm

import (
	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/log"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/infrastructure"
)

type MapReducer struct {
	logger   log.ILogger
	sharding infrastructure.Sharding
	model    entity.Entity
}

var _ entity.MapReducer = (*MapReducer)(nil)

func NewMapReducer(logger log.ILogger, sharding infrastructure.Sharding) *MapReducer {
	return &MapReducer{
		logger:   logger,
		sharding: sharding,
		model:    entity.Entity{},
	}
}

func (s *MapReducer) GetDBByID(ID uint, typeID uint) minipkg_gorm.IDB {
	return nil
}

func (s *MapReducer) GetDBForNew(typeID uint) minipkg_gorm.IDB {
	return nil
}
