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

func (s *MapReducer) GetDB(ctx context.Context, typeID uint, ID uint) (minipkg_gorm.IDB, error) {
	entityType, err := s.entityTypeFinder.GetSysnameByID(ctx, typeID)
	if err != nil {
		return nil, err
	}

	cluster := s.sharding.Default
	if sysname, ok := s.sharding.ClusterSysnamesByEntityTypes[entityType]; ok {
		if cluster, ok = s.sharding.BySysnames[sysname]; !ok {
			return nil, errors.Wrapf(yaruserror.ErrNotFound, "Cluster not found, sysname = %q.", sysname)
		}
	}

	shardIndex := s.ShardIndex(cluster.Capacity, ID)
	if shardIndex >= uint(len(cluster.Items)) {
		return nil, errors.Errorf("ID = %v is too big for cluster capacity = %v length = %v", ID, cluster.Capacity, len(cluster.Items))
	}

	return cluster.Items[shardIndex], nil
}

func (s *MapReducer) GetDBs(parser *SelectionConditionParser, entityWhereConditions selection_condition.WhereConditions) ([]minipkg_gorm.IDB, error) {
	//	получаем слайс EntityType и слайс ID
	entityTypes, IDs, err := parser.GetEntityTypeIDsAndIDsByEntityWhereConditions(entityWhereConditions)
	if err != nil {
		return nil, err
	}
	//	получаем слайс кластеров по слайсу EntityType
	//	бежим по слайсу кластеров:
	//		получаем слайс db по слайсу ID
	return s.sharding.GetDBs(condition)
}

func (s *MapReducer) Query(ctx context.Context, parser *SelectionConditionParser, entityWhereConditions selection_condition.WhereConditions, f func(db minipkg_gorm.IDB) ([]SearchResult, error)) ([]SearchResult, error) {
	var res []SearchResult

	dbs, err := s.GetDBs(parser, entityWhereConditions)
	if err != nil {
		return nil, err
	}

	for _, db := range dbs {
		searchResult, err := f(db)
		// todo: распараллелить
		if err != nil && !errors.Is(err, yaruserror.ErrNotFound) {
			return nil, err
		}
		res = append(res, searchResult...)
	}
	// todo: сортировку результатов
	return res, nil
}

func (s *MapReducer) Count(ctx context.Context, parser *SelectionConditionParser, entityWhereConditions selection_condition.WhereConditions, f func(db minipkg_gorm.IDB) (uint, error)) (uint, error) {
	var res uint

	dbs, err := s.GetDBs(parser, entityWhereConditions)
	if err != nil {
		return 0, err
	}

	for _, db := range dbs {
		searchResult, err := f(db)
		if err != nil && !errors.Is(err, yaruserror.ErrNotFound) {
			return 0, err
		}
		res += searchResult
	}

	return res, nil
}
