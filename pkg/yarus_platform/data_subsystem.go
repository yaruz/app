package yarus_platform

import (
	golog "log"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/property_type"

	"github.com/minipkg/log"
	"github.com/pkg/errors"
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
	Service    property_type.IService
	Repository property_type.Repository
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

func newDataDomain(infra *infrastructure) (*DataDomain, error) {
	d := &DataDomain{}
	if err := d.setupRepositories(infra); err != nil {
		return nil, err
	}
	d.setupServices(infra.Logger)
	return d, nil
}

func (d *DataDomain) setupRepositories(infra *infrastructure) (err error) {
	var ok bool

	repo, err := gormrep.GetRepository(infra.Logger, infra.DataDB, t_source.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", t_source.EntityName, err)
	}
	d.TSource.Repository, ok = repo.(t_source.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", t_source.EntityName, t_source.EntityName, repo)
	}

	repo, err = gormrep.GetRepository(infra.Logger, infra.DataDB, property_unit.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_unit.EntityName, err)
	}
	d.PropertyUnit.Repository, ok = repo.(property_unit.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_unit.EntityName, property_unit.EntityName, repo)
	}

	repo, err = gormrep.GetRepository(infra.Logger, infra.DataDB, property_group.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_group.EntityName, err)
	}
	d.PropertyGroup.Repository, ok = repo.(property_group.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_group.EntityName, property_group.EntityName, repo)
	}

	repo, err = gormrep.GetRepository(infra.Logger, infra.DataDB, property_type.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_type.EntityName, err)
	}
	d.PropertyType.Repository, ok = repo.(property_type.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_type.EntityName, property_type.EntityName, repo)
	}

	repo, err = gormrep.GetRepository(infra.Logger, infra.DataDB, property_view_type.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_view_type.EntityName, err)
	}
	d.PropertyViewType.Repository, ok = repo.(property_view_type.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_view_type.EntityName, property_view_type.EntityName, repo)
	}

	repo, err = gormrep.GetRepository(infra.Logger, infra.DataDB, property_type2property_view_type.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_type2property_view_type.EntityName, err)
	}
	d.PropertyType2PropertyViewType.Repository, ok = repo.(property_type2property_view_type.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_type2property_view_type.EntityName, property_type2property_view_type.EntityName, repo)
	}

	repo, err = gormrep.GetRepository(infra.Logger, infra.DataDB, property.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property.EntityName, err)
	}
	d.Property.Repository, ok = repo.(property.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property.EntityName, property.EntityName, repo)
	}

	repo, err = gormrep.GetRepository(infra.Logger, infra.DataDB, t_string.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", t_string.EntityName, err)
	}
	d.TString.Repository, ok = repo.(t_string.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", t_string.EntityName, t_string.EntityName, repo)
	}

	repo, err = gormrep.GetRepository(infra.Logger, infra.DataDB, t_text.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", t_text.EntityName, err)
	}
	d.TText.Repository, ok = repo.(t_text.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", t_text.EntityName, t_text.EntityName, repo)
	}

	repo, err = gormrep.GetRepository(infra.Logger, infra.DataDB, entity_type.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", entity_type.EntityName, err)
	}
	d.EntityType.Repository, ok = repo.(entity_type.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", entity_type.EntityName, entity_type.EntityName, repo)
	}

	repo, err = gormrep.GetRepository(infra.Logger, infra.DataDB, entity_type2property.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", entity_type2property.EntityName, err)
	}
	d.EntityType2Property.Repository, ok = repo.(entity_type2property.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", entity_type2property.EntityName, entity_type2property.EntityName, repo)
	}

	repo, err = gormrep.GetRepository(infra.Logger, infra.DataDB, entity.EntityName)
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
	d.Entity.Service = entity.NewService(logger, d.Entity.Repository)
	d.EntityType2Property.Service = entity_type2property.NewService(logger, d.EntityType2Property.Repository)
	d.EntityType.Service = entity_type.NewService(logger, d.EntityType.Repository)
	d.PropertyGroup.Service = property_group.NewService(logger, d.PropertyGroup.Repository)
	d.Property.Service = property.NewService(logger, d.Property.Repository)
	d.PropertyType2PropertyViewType.Service = property_type2property_view_type.NewService(logger, d.PropertyType2PropertyViewType.Repository)
	d.PropertyType.Service = property_type.NewService(logger, d.PropertyType.Repository)
	d.PropertyUnit.Service = property_unit.NewService(logger, d.PropertyUnit.Repository)
	d.PropertyViewType.Service = property_view_type.NewService(logger, d.PropertyViewType.Repository)
	d.TSource.Service = t_source.NewService(logger, d.TSource.Repository)
	d.TString.Service = t_string.NewService(logger, d.TString.Repository)
	d.TText.Service = t_text.NewService(logger, d.TText.Repository)
}
