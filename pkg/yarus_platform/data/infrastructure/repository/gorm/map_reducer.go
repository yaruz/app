package gorm

import (
	"context"

	"github.com/pkg/errors"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/log"
	"github.com/minipkg/selection_condition"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/infrastructure"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
)

type MapReducer struct {
	logger           log.ILogger
	sharding         infrastructure.Sharding
	model            entity.Entity
	entityTypeFinder entity.EntityTypeFinder
}

var _ IMapReducer = (*MapReducer)(nil)

func NewMapReducer(logger log.ILogger, entityTypeFinder entity.EntityTypeFinder, sharding infrastructure.Sharding) *MapReducer {
	return &MapReducer{
		logger:           logger,
		entityTypeFinder: entityTypeFinder,
		sharding:         sharding,
		model:            entity.Entity{},
	}
}

func (s *MapReducer) ShardIndex(shardCapacity uint, ID uint) uint {
	res := ID / shardCapacity
	ost := ID % shardCapacity

	if ost > 0 {
		res++
	}
	return res - 1
}

func (s *MapReducer) GetDB(ctx context.Context, ID uint, typeID uint) (minipkg_gorm.IDB, error) {
	typeSysname, err := s.entityTypeFinder.GetSysnameByID(ctx, typeID)
	if err != nil {
		return nil, err
	}

	shards := &s.sharding.Default
	if s, ok := s.sharding.ByTypes[typeSysname]; ok {
		shards = &s
	}

	shardIndex := s.ShardIndex(shards.Capacity, ID)
	if shardIndex >= uint(len(shards.Items)) {
		return nil, errors.Errorf("ID = %v is too big for shards capacity = %v length = %v", ID, shards.Capacity, len(shards.Items))
	}

	return (*shards).Items[shardIndex], nil
}

//	1. централизованные последовательности sequence для каждого списка шард (Shards). Храним в DB reference.
//	2. соотв. добываем и инкрементим его программно. И пишем в entity.ID - простое поле, не sequence.
//	3. Соотв. для этого нужна соотв. структура со списком методов

func (s *MapReducer) GetDBs(condition *selection_condition.SelectionCondition) []minipkg_gorm.IDB {
	return s.sharding.GetDBs(condition)
}

func (s *MapReducer) GetDBForInsert(typeID uint) minipkg_gorm.IDB {
	return s.sharding.GetDBForInsert(typeID)
}

func (s *MapReducer) Query(ctx context.Context, model interface{}, condition *selection_condition.SelectionCondition, f func(db minipkg_gorm.IDB) ([]SearchResult, error)) ([]SearchResult, error) {
	var res []SearchResult

	dbs := s.GetDBs(condition)
	for _, db := range dbs {
		searchResult, err := f(db)
		if err != nil && !errors.Is(err, yaruserror.ErrNotFound) {
			return nil, err
		}
		res = append(res, searchResult...)
	}
	// todo: сортировку результатов
	return res, nil
}

func (s *MapReducer) Count(ctx context.Context, model interface{}, condition *selection_condition.SelectionCondition, f func(db minipkg_gorm.IDB) (uint, error)) (uint, error) {
	var res uint

	dbs := s.GetDBs(condition)
	for _, db := range dbs {
		searchResult, err := f(db)
		if err != nil && !errors.Is(err, yaruserror.ErrNotFound) {
			return 0, err
		}
		res += searchResult
	}

	return res, nil
}
