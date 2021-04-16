package gorm

import (
	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity_type"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity_type2property"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/property_group"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/property_type"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/property_type2property_view_type"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/property_unit"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/property_view_type"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_source"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/minipkg/log"
)

// IRepository is an interface of repository
type IRepository interface{}

// repository persists albums in database
type repository struct {
	db         minipkg_gorm.IDB
	logger     log.ILogger
	Conditions *selection_condition.SelectionCondition
}

const DefaultLimit = 1000

// GetRepository return a repository
func GetRepository(logger log.ILogger, dbase minipkg_gorm.IDB, entityName string) (repo IRepository, err error) {
	r := &repository{
		db:     dbase,
		logger: logger,
	}

	switch entityName {
	case entity.EntityName:
		repo, err = NewEntityRepository(r)
	case entity_type2property.EntityName:
		repo, err = NewEntityType2PropertyRepository(r)
	case entity_type.EntityName:
		repo, err = NewEntityTypeRepository(r)
	case property_group.EntityName:
		repo, err = NewPropertyGroupRepository(r)
	case property.EntityName:
		repo, err = NewPropertyRepository(r)
	case property_type2property_view_type.EntityName:
		repo, err = NewPropertyType2PropertyViewTypeRepository(r)
	case property_type.EntityName:
		repo, err = NewPropertyTypeRepository(r)
	case property_unit.EntityName:
		repo, err = NewPropertyUnitRepository(r)
	case property_view_type.EntityName:
		repo, err = NewPropertyViewTypeRepository(r)
	case text_source.EntityName:
		repo, err = NewTextSourceRepository(r)
	case text_value.EntityName:
		repo, err = NewTextValueRepository(r)
	default:
		err = errors.Errorf("Repository for entity %q not found", entityName)
	}
	return repo, err
}

func (r *repository) SetDefaultConditions(defaultConditions *selection_condition.SelectionCondition) {
	r.Conditions = defaultConditions
}

func (r repository) DB() *gorm.DB {
	return minipkg_gorm.Conditions(r.db.DB(), r.Conditions)
}
