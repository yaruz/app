package yarus_platform

import (
	golog "log"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/db/redis"
	"github.com/minipkg/db/redis/cache"
	"github.com/minipkg/log"
	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/config"
	gormrep "github.com/yaruz/app/pkg/yarus_platform/data/infrastructure/repository/gorm"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity_type"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity_type2property"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/property_group"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/property_type2property_view_type"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/property_unit"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/property_view_type"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/t_source"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/t_string"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/t_text"
)

type infrastructure struct {
	DataDB   minipkg_gorm.IDB
	SearchDB minipkg_gorm.IDB
	Redis    redis.IDB
	Cache    cache.Service
}

func newInfra(logger log.ILogger, cfg config.Infrastructure) (*infrastructure, error) {

	DataDB, err := minipkg_gorm.New(cfg.DataDB, logger)
	if err != nil {
		return nil, err
	}

	SearchDB, err := minipkg_gorm.New(cfg.SearchDB, logger)
	if err != nil {
		return nil, err
	}

	rDB, err := redis.New(cfg.Redis)
	if err != nil {
		return nil, err
	}

	return &infrastructure{
		DataDB:   DataDB,
		SearchDB: SearchDB,
		Redis:    rDB,
		Cache:    cache.NewService(rDB, cfg.CacheLifeTime),
	}, nil
}

type DataDomain struct {
	Entity                        DataDomainEntity
	EntityType                    DataDomainEntityType
	EntityType2Property           DataDomainEntityType2Property
	Property                      DataDomainProperty
	PropertyGroup                 DataDomainPropertyGroup
	PropertyType                  DataDomainPropertyType
	PropertyType2PropertyViewType DataDomainPropertyType2PropertyViewType
	PropertyUnit                  DataDomainPropertyUnit
	PropertyViewType              DataDomainPropertyViewType
	TSource                       DataDomainTSource
	TString                       DataDomainTString
	TText                         DataDomainTText
}

func newDataDomain(logger log.ILogger, infra *infrastructure) (*DataDomain, error) {
	d := &DataDomain{}
	if err := d.setupRepositories(logger, infra); err != nil {
		return nil, err
	}
	d.setupServices(logger)
	return d, nil
}

func (d *DataDomain) setupRepositories(logger log.ILogger, infra *infrastructure) (err error) {
	var ok bool

	repo, err := gormrep.GetRepository(logger, infra.DataDB, entity.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", entity.EntityName, err)
	}

	d.Entity.Repository, ok = repo.(entity.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", entity.EntityName, entity.EntityName, repo)
	}

	return nil
}

func (d *DataDomain) setupServices(logger log.ILogger) {

}

type SearchDomain struct {
}

func newSearchDomain(logger log.ILogger, infra *infrastructure) (*SearchDomain, error) {
	d := &SearchDomain{}
	if err := d.setupRepositories(logger, infra); err != nil {
		return nil, err
	}
	d.setupServices(logger)
	return d, nil
}

func (d *SearchDomain) setupRepositories(logger log.ILogger, infra *infrastructure) (err error) {
	return nil
}

func (d *SearchDomain) setupServices(logger log.ILogger) {

}

type DataDomainEntity struct {
	Service    entity.IService
	Repository entity.Repository
}

type DataDomainEntityType struct {
	Service    entity_type.IService
	Repository entity_type.Repository
}

type DataDomainEntityType2Property struct {
	Service    entity_type2property.IService
	Repository entity_type2property.Repository
}

type DataDomainProperty struct {
	Service    property.IService
	Repository property.Repository
}

type DataDomainPropertyGroup struct {
	Service    property_group.IService
	Repository property_group.Repository
}

type DataDomainPropertyType struct {
	Service    property_group.IService
	Repository property_group.Repository
}

type DataDomainPropertyType2PropertyViewType struct {
	Service    property_type2property_view_type.IService
	Repository property_type2property_view_type.Repository
}

type DataDomainPropertyUnit struct {
	Service    property_unit.IService
	Repository property_unit.Repository
}

type DataDomainPropertyViewType struct {
	Service    property_view_type.IService
	Repository property_view_type.Repository
}

type DataDomainTSource struct {
	Service    t_source.IService
	Repository t_source.Repository
}

type DataDomainTString struct {
	Service    t_string.IService
	Repository t_string.Repository
}

type DataDomainTText struct {
	Service    t_text.IService
	Repository t_text.Repository
}
