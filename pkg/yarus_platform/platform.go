package yarus_platform

import (
	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/db/redis"
	"github.com/minipkg/db/redis/cache"
	"github.com/minipkg/log"
	"github.com/yaruz/app/pkg/yarus_platform/config"

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

func newInfra(cfg config.Infrastructure, logger log.ILogger) (*infrastructure, error) {

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

func newDataDomain() (*DataDomain, error) {

}

type SearchDomain struct {
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
