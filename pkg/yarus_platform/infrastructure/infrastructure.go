package infrastructure

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/db/redis"
	"github.com/minipkg/db/redis/cache"
	"github.com/minipkg/log"
	"github.com/pkg/errors"
	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/pkg/yarus_platform/config"
)

type Infrastructure struct {
	Logger       log.ILogger
	DataSharding Sharding
	ReferenceDB  minipkg_gorm.IDB
	SearchDB     minipkg_gorm.IDB
	Redis        redis.IDB
	Cache        cache.Service
}

type Sharding struct {
	IsAutoMigrate bool
	Model         interface{}
	Default       Shards
	ByTypes       map[string]Shards
}

func (s *Sharding) SchemesInitWithContext(ctx context.Context) (err error) {
	return s.ApplyFunc2DBs(func(db minipkg_gorm.IDB) (err error) {
		db, err = db.SchemeInitWithContext(ctx, s.Model)
		return err
	})
}

func (s *Sharding) Close() (err error) {
	return s.ApplyFunc2DBs(func(db minipkg_gorm.IDB) error {
		return s.ApplyFunc2DBs(func(db minipkg_gorm.IDB) error {
			return db.Close()
		})
	})
}

func (s *Sharding) ApplyFunc2DBs(f func(db minipkg_gorm.IDB) error) (err error) {
	for _, shards := range s.ByTypes {
		if err = shards.ApplyFunc2DBs(f); err != nil {
			return err
		}
	}
	return s.Default.ApplyFunc2DBs(f)
}

type Shards struct {
	Capacity uint
	Items    []minipkg_gorm.IDB
}

func (s *Shards) SchemesInitWithContext(ctx context.Context, model interface{}) (err error) {
	return s.ApplyFunc2DBs(func(db minipkg_gorm.IDB) (err error) {
		db, err = db.SchemeInitWithContext(ctx, model)
		return err
	})
}

func (s *Shards) Close() (err error) {
	return s.ApplyFunc2DBs(func(db minipkg_gorm.IDB) error {
		return db.Close()
	})
}

func (s *Shards) ApplyFunc2DBs(f func(db minipkg_gorm.IDB) error) (err error) {
	for i := range s.Items {
		if err = f(s.Items[i]); err != nil {
			return err
		}
	}
	return nil
}

func NewInfrastructure(ctx context.Context, logger log.ILogger, cfg *config.Infrastructure) (*Infrastructure, error) {

	DataSharding, err := newDataSharding(ctx, logger, &cfg.DataSharding, entity.New())
	if err != nil {
		return nil, err
	}

	ReferenceDB, err := minipkg_gorm.New(logger, cfg.ReferenceDB)
	if err != nil {
		return nil, err
	}

	SearchDB, err := minipkg_gorm.New(logger, cfg.SearchDB)
	if err != nil {
		return nil, err
	}

	rDB, err := redis.New(cfg.Redis)
	if err != nil {
		return nil, err
	}

	return &Infrastructure{
		Logger:       logger,
		DataSharding: *DataSharding,
		ReferenceDB:  ReferenceDB,
		SearchDB:     SearchDB,
		Redis:        rDB,
		Cache:        cache.NewService(rDB, cfg.CacheLifeTime),
	}, nil
}

func newDataSharding(ctx context.Context, logger log.ILogger, cfg *config.Sharding, model interface{}) (*Sharding, error) {
	defaultShards, err := newShards(logger, &cfg.Default)
	if err != nil {
		return nil, err
	}

	byTypes := make(map[string]Shards, len(cfg.ByTypes))
	for t, shardsCfg := range cfg.ByTypes {
		shards, err := newShards(logger, &shardsCfg)
		if err != nil {
			return nil, err
		}
		byTypes[t] = *shards
	}

	s := &Sharding{
		IsAutoMigrate: cfg.IsAutoMigrate,
		Model:         model,
		Default:       *defaultShards,
		ByTypes:       byTypes,
	}

	if err := s.SchemesInitWithContext(ctx); err != nil {
		return nil, err
	}
	return s, nil
}

func newShards(logger log.ILogger, cfg *config.Shards) (*Shards, error) {
	var err error

	items := make([]minipkg_gorm.IDB, len(cfg.Items))
	for i, cfgItem := range cfg.Items {
		if items[i], err = minipkg_gorm.New(logger, cfgItem); err != nil {
			return nil, err
		}
	}

	return &Shards{
		Capacity: cfg.Capacity,
		Items:    items,
	}, nil
}

func (i *Infrastructure) Stop() error {
	errRedis := i.Redis.Close()
	errDataDB := i.DataSharding.Close()
	errReferenceDB := i.ReferenceDB.Close()
	errSearchDB := i.SearchDB.Close()

	switch {
	case errDataDB != nil:
		return errors.Wrapf(apperror.ErrInternal, "db close error: %v", errDataDB)
	case errReferenceDB != nil:
		return errors.Wrapf(apperror.ErrInternal, "db close error: %v", errReferenceDB)
	case errSearchDB != nil:
		return errors.Wrapf(apperror.ErrInternal, "db close error: %v", errSearchDB)
	case errRedis != nil:
		return errors.Wrapf(apperror.ErrInternal, "redis close error: %v", errRedis)
	}

	return nil
}
