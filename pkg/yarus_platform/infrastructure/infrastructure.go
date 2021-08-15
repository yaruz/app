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
	Logger      log.ILogger
	DataDB      minipkg_gorm.IDB
	ReferenceDB minipkg_gorm.IDB
	SearchDB    minipkg_gorm.IDB
	Redis       redis.IDB
	Cache       cache.Service
}

func NewInfra(logger log.ILogger, cfg config.Infrastructure) (*Infrastructure, error) {

	DataDB, err := minipkg_gorm.New(logger, cfg.DataDB)
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
		Logger:      logger,
		DataDB:      DataDB,
		ReferenceDB: ReferenceDB,
		SearchDB:    SearchDB,
		Redis:       rDB,
		Cache:       cache.NewService(rDB, cfg.CacheLifeTime),
	}, nil
}

func (i *Infrastructure) Stop() error {
	errRedis := i.Redis.Close()
	errDataDB := i.DataDB.Close()
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
