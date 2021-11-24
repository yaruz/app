package infrastructure

import (
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
	Default Shards
	ByTypes map[string]Shards
}

func (s *Sharding) Close() (err error) {
	for _, shards := range s.ByTypes {
		if err = shards.Close(); err != nil {
			return err
		}
	}
	return s.Default.Close()
}

type Shards struct {
	Capacity uint
	Items    []minipkg_gorm.IDB
}

func (s *Shards) Close() (err error) {
	for i := range s.Items {
		if err = s.Items[i].Close(); err != nil {
			return err
		}
	}
	return nil
}

func NewInfrastructure(logger log.ILogger, cfg *config.Infrastructure) (*Infrastructure, error) {

	DataSharding, err := newDataSharding(logger, &cfg.DataSharding)
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

func newDataSharding(logger log.ILogger, cfg *config.Sharding) (*Sharding, error) {
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

	return &Sharding{
		Default: *defaultShards,
		ByTypes: byTypes,
	}, nil
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
