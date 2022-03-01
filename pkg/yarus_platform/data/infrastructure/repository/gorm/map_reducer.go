package gorm

import (
	"context"
	"fmt"
	"sync"

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

func (s *MapReducer) GetDBs(ctx context.Context, parser *SelectionConditionParser, entityWhereConditions selection_condition.WhereConditions) ([]minipkg_gorm.IDB, error) {
	var dbs []minipkg_gorm.IDB
	//	получаем слайс EntityType и слайс ID
	entityTypeIDs, IDs, err := parser.GetEntityTypeIDsAndIDsByEntityWhereConditions(entityWhereConditions)
	if err != nil {
		return nil, err
	}

	entityTypeSysnames := make([]string, len(entityTypeIDs))
	for i, entityTypeID := range entityTypeIDs {
		entityTypeSysname, err := s.entityTypeFinder.GetSysnameByID(ctx, entityTypeID)
		if err != nil {
			return nil, err
		}
		entityTypeSysnames[i] = entityTypeSysname
	}

	//	получаем слайс кластеров по слайсу EntityType
	clusters, err := s.sharding.GetClustersByEntityTypes(entityTypeSysnames)
	if err != nil {
		return nil, err
	}

	//	бежим по слайсу кластеров:
	for _, cluster := range clusters {
		if IDs == nil {
			dbs = append(dbs, cluster.Items...)
			continue
		}
		//		получаем слайс db по слайсу ID
		for _, ID := range IDs {
			shardIndex := s.ShardIndex(cluster.Capacity, ID)
			if shardIndex >= uint(len(cluster.Items)) {
				return nil, errors.Errorf("ID = %v is too big for cluster capacity = %v length = %v", ID, cluster.Capacity, len(cluster.Items))
			}
			dbs = append(dbs, cluster.Items[shardIndex])
		}

	}

	return dbs, nil
}

func (s *MapReducer) Query(ctx context.Context, parser *SelectionConditionParser, entityWhereConditions selection_condition.WhereConditions, f func(db minipkg_gorm.IDB) ([]SearchResult, error)) ([]SearchResult, error) {
	dbs, err := s.GetDBs(ctx, parser, entityWhereConditions)
	if err != nil {
		return nil, err
	}

	// todo: сортировку результатов
	return s.queryReceiver(s.queryStarter(dbs, f))
}

func (s *MapReducer) queryStarter(dbs []minipkg_gorm.IDB, f func(db minipkg_gorm.IDB) ([]SearchResult, error)) (chan []SearchResult, chan error) {
	var errorsCh = make(chan error, len(dbs))
	var resultCh = make(chan []SearchResult, len(dbs))

	select {
	case err, ok := <-errorsCh:
		fmt.Println("Какого хуя?!", err, ok)
	default:
		fmt.Println("Хуй!")
	}

	defer func() {
		close(errorsCh)
		close(resultCh)
	}()

	wg := &sync.WaitGroup{}
	for _, db := range dbs {
		wg.Add(1)
		go s.queryProcessing(wg, db, f, resultCh, errorsCh)
	}
	wg.Wait()

	select {
	case err := <-errorsCh:
		fmt.Println("Какого хуя?!", err)
	default:
		fmt.Println("Хуй!")
	}

	return resultCh, errorsCh
}

func (s *MapReducer) queryProcessing(wg *sync.WaitGroup, db minipkg_gorm.IDB, f func(db minipkg_gorm.IDB) ([]SearchResult, error), resultCh chan []SearchResult, errorsCh chan error) {
	var searchResult []SearchResult
	defer func() {
		if err := recover(); err != nil {
			errorsCh <- errors.New("panic happened in given func: " + err.(string))
		}
		wg.Done()
	}()

	searchResult, err := f(db)
	if err != nil {
		if !errors.Is(err, yaruserror.ErrNotFound) {
			errorsCh <- err
		}
		return
	}

	resultCh <- searchResult
	return
}

func (s *MapReducer) queryReceiver(resultCh chan []SearchResult, errorsCh chan error) ([]SearchResult, error) {
	var res []SearchResult
	var err error

	select {
	case err = <-errorsCh:
		return nil, err
	default:
	}

	for searchResult := range resultCh {
		res = append(res, searchResult...)
	}

	if res == nil {
		err = yaruserror.ErrNotFound
	}
	return res, err
}

func (s *MapReducer) Count(ctx context.Context, parser *SelectionConditionParser, entityWhereConditions selection_condition.WhereConditions, f func(db minipkg_gorm.IDB) (uint, error)) (uint, error) {
	dbs, err := s.GetDBs(ctx, parser, entityWhereConditions)
	if err != nil {
		return 0, err
	}

	return s.countReceiver(s.countStarter(dbs, f))
}

func (s *MapReducer) countStarter(dbs []minipkg_gorm.IDB, f func(db minipkg_gorm.IDB) (uint, error)) (resultCh chan uint, errorsCh chan error) {
	errorsCh = make(chan error)
	resultCh = make(chan uint)

	defer func() {
		close(errorsCh)
		close(resultCh)
	}()

	wg := &sync.WaitGroup{}
	for _, db := range dbs {
		wg.Add(1)
		go s.countProcessing(wg, db, f, resultCh, errorsCh)
	}
	wg.Wait()

	return resultCh, errorsCh
}

func (s *MapReducer) countProcessing(wg *sync.WaitGroup, db minipkg_gorm.IDB, f func(db minipkg_gorm.IDB) (uint, error), resultCh chan uint, errorsCh chan error) {
	var searchResult uint
	defer func() {
		if err := recover(); err != nil {
			errorsCh <- errors.New("panic happened in given func: " + err.(string))
		}
		wg.Done()
	}()

	searchResult, err := f(db)
	if err != nil {
		if !errors.Is(err, yaruserror.ErrNotFound) {
			errorsCh <- err
		}
		return
	}

	resultCh <- searchResult
	return
}

func (s *MapReducer) countReceiver(resultCh chan uint, errorsCh chan error) (uint, error) {
	var res uint
	var err error

	select {
	case err = <-errorsCh:
		return 0, err
	default:
	}

	for searchResult := range resultCh {
		res += searchResult
	}
	return res, nil
}
