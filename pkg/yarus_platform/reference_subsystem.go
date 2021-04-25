package yarus_platform

import (
	golog "log"

	"github.com/yaruz/app/pkg/yarus_platform/reference/infrastructure/repository/gorm"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"

	"github.com/minipkg/log"
	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type2property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_group"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type2property_view_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_unit"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_view_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_source"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_value"
)

type ReferenceDomain struct {
	EntityType                    ReferenceDomainEntityType
	EntityType2Property           ReferenceDomainEntityType2Property
	Property                      ReferenceDomainProperty
	PropertyGroup                 ReferenceDomainPropertyGroup
	PropertyType                  ReferenceDomainPropertyType
	PropertyType2PropertyViewType ReferenceDomainPropertyType2PropertyViewType
	PropertyUnit                  ReferenceDomainPropertyUnit
	PropertyViewType              ReferenceDomainPropertyViewType
	TextSource                    ReferenceDomainTextSource
	TextValue                     ReferenceDomainTextValue
}

type ReferenceDomainEntityType struct {
	Service    entity_type.IService
	Repository entity_type.Repository
}

type ReferenceDomainEntityType2Property struct {
	Service    entity_type2property.IService
	Repository entity_type2property.Repository
}

type ReferenceDomainProperty struct {
	Service    property.IService
	Repository property.Repository
}

type ReferenceDomainPropertyGroup struct {
	Service    property_group.IService
	Repository property_group.Repository
}

type ReferenceDomainPropertyType struct {
	Service    property_type.IService
	Repository property_type.Repository
}

type ReferenceDomainPropertyType2PropertyViewType struct {
	Service    property_type2property_view_type.IService
	Repository property_type2property_view_type.Repository
}

type ReferenceDomainPropertyUnit struct {
	Service    property_unit.IService
	Repository property_unit.Repository
}

type ReferenceDomainPropertyViewType struct {
	Service    property_view_type.IService
	Repository property_view_type.Repository
}

type ReferenceDomainTextSource struct {
	Service    text_source.IService
	Repository text_source.Repository
}

type ReferenceDomainTextValue struct {
	Service    text_value.IService
	Repository text_value.Repository
}

func newReferenceDomain(infra *infrastructure) (*ReferenceDomain, error) {
	d := &ReferenceDomain{}
	if err := d.setupRepositories(infra); err != nil {
		return nil, err
	}
	d.setupServices(infra.Logger)
	return d, nil
}

func (d *ReferenceDomain) setupRepositories(infra *infrastructure) (err error) {
	var ok bool

	repo, err := gorm.GetRepository(infra.Logger, infra.ReferenceDB, text_source.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", text_source.EntityName, err)
	}
	d.TextSource.Repository, ok = repo.(text_source.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", text_source.EntityName, text_source.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property_unit.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_unit.EntityName, err)
	}
	d.PropertyUnit.Repository, ok = repo.(property_unit.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_unit.EntityName, property_unit.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property_group.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_group.EntityName, err)
	}
	d.PropertyGroup.Repository, ok = repo.(property_group.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_group.EntityName, property_group.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property_view_type.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_view_type.EntityName, err)
	}
	d.PropertyViewType.Repository, ok = repo.(property_view_type.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_view_type.EntityName, property_view_type.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property_type.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_type.EntityName, err)
	}
	d.PropertyType.Repository, ok = repo.(property_type.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_type.EntityName, property_type.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property_type2property_view_type.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property_type2property_view_type.EntityName, err)
	}
	d.PropertyType2PropertyViewType.Repository, ok = repo.(property_type2property_view_type.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property_type2property_view_type.EntityName, property_type2property_view_type.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, property.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", property.EntityName, err)
	}
	d.Property.Repository, ok = repo.(property.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property.EntityName, property.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, text_value.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", text_value.EntityName, err)
	}
	d.TextValue.Repository, ok = repo.(text_value.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", text_value.EntityName, text_value.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, entity_type.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", entity_type.EntityName, err)
	}
	d.EntityType.Repository, ok = repo.(entity_type.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", entity_type.EntityName, entity_type.EntityName, repo)
	}

	repo, err = gorm.GetRepository(infra.Logger, infra.ReferenceDB, entity_type2property.EntityName)
	if err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", entity_type2property.EntityName, err)
	}
	d.EntityType2Property.Repository, ok = repo.(entity_type2property.Repository)
	if !ok {
		return errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", entity_type2property.EntityName, entity_type2property.EntityName, repo)
	}

	return nil
}

func (d *ReferenceDomain) setupServices(logger log.ILogger) {
	d.EntityType2Property.Service = entity_type2property.NewService(logger, d.EntityType2Property.Repository)
	d.EntityType.Service = entity_type.NewService(logger, d.EntityType.Repository)
	d.PropertyGroup.Service = property_group.NewService(logger, d.PropertyGroup.Repository)
	d.Property.Service = property.NewService(logger, d.Property.Repository)
	d.PropertyType2PropertyViewType.Service = property_type2property_view_type.NewService(logger, d.PropertyType2PropertyViewType.Repository)
	d.PropertyType.Service = property_type.NewService(logger, d.PropertyType.Repository)
	d.PropertyUnit.Service = property_unit.NewService(logger, d.PropertyUnit.Repository)
	d.PropertyViewType.Service = property_view_type.NewService(logger, d.PropertyViewType.Repository)
	d.TextSource.Service = text_source.NewService(logger, d.TextSource.Repository)
	d.TextValue.Service = text_value.NewService(logger, d.TextValue.Repository)
}
