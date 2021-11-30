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
	logger   log.ILogger
	sharding infrastructure.Sharding
	model    entity.Entity
}

var _ IMapReducer = (*MapReducer)(nil)

func NewMapReducer(logger log.ILogger, sharding infrastructure.Sharding) *MapReducer {
	return &MapReducer{
		logger:   logger,
		sharding: sharding,
		model:    entity.Entity{},
	}
}

func (s *MapReducer) GetDB(ID uint, typeID uint) minipkg_gorm.IDB {
	return s.sharding.GetDB(ID, typeID)
}

func (s *MapReducer) GetDBs(condition *selection_condition.SelectionCondition) []minipkg_gorm.IDB {
	return s.sharding.GetDBs(condition)
}

func (s *MapReducer) GetDBForInsert(typeID uint) minipkg_gorm.IDB {
	return s.sharding.GetDBForInsert(typeID)
}

func (s *MapReducer) Query(ctx context.Context, model interface{}, condition *selection_condition.SelectionCondition, f func(db minipkg_gorm.IDB) ([]SearchResult, error)) ([]SearchResult, error) {
	var res []SearchResult

	dbs := s.sharding.GetDBs(condition)
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

	dbs := s.sharding.GetDBs(condition)
	for _, db := range dbs {
		searchResult, err := f(db)
		if err != nil && !errors.Is(err, yaruserror.ErrNotFound) {
			return 0, err
		}
		res += searchResult
	}

	return res, nil
}
