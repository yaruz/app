package gorm

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/relation"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/log"
	"github.com/minipkg/selection_condition"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type2property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_group"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type2property_view_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_unit"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_view_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_source"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_value"
)

// IRepository is an interface of repository
type IRepository interface {
	DB() *gorm.DB
}

// repository persists albums in database
type repository struct {
	db         minipkg_gorm.IDB
	logger     log.ILogger
	Conditions *selection_condition.SelectionCondition
	model      interface{}
}

const DefaultLimit = 1000

// GetRepository return a repository
func GetRepository(logger log.ILogger, dbase minipkg_gorm.IDB, entityName string) (repo IRepository, err error) {
	r := &repository{
		logger: logger,
	}
	ctx := context.Background()

	switch entityName {
	case entity_type2property.EntityName:
		r.model = entity_type2property.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewEntityType2PropertyRepository(r)
	case entity_type.EntityName:
		entityType2PropertyRepo, err := GetRepository(logger, dbase, entity_type2property.EntityName)
		if err != nil {
			return nil, errors.Wrapf(err, "Can not get db repository for entity %q, error happened: %v", entity_type2property.EntityName, err)
		}
		entityType2PropertyRepository, ok := entityType2PropertyRepo.(entity_type2property.Repository)
		if !ok {
			return nil, errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", entity_type2property.EntityName, entity_type2property.EntityName, entityType2PropertyRepo)
		}

		r.model = entity_type.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewEntityTypeRepository(r, &entityType2PropertyRepository)
	case property_group.EntityName:
		r.model = property_group.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewPropertyGroupRepository(r)
	case property.EntityName:
		r.model = property.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewPropertyRepository(r)
	case property_type2property_view_type.EntityName:
		r.model = property_type2property_view_type.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewPropertyType2PropertyViewTypeRepository(r)
	case property_type.EntityName:
		r.model = property_type.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewPropertyTypeRepository(r)
	case property_unit.EntityName:
		r.model = property_unit.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewPropertyUnitRepository(r)
	case property_view_type.EntityName:
		r.model = property_view_type.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewPropertyViewTypeRepository(r)
	case text_source.EntityName:
		r.model = text_source.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewTextSourceRepository(r)
	case text_value.EntityName:
		r.model = text_value.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewTextValueRepository(r)
	case relation.EntityName:
		entityType2PropertyRepo, err := GetRepository(logger, dbase, entity_type2property.EntityName)
		if err != nil {
			return nil, errors.Wrapf(err, "Can not get db repository for entity %q, error happened: %v", entity_type2property.EntityName, err)
		}
		entityType2PropertyRepository, ok := entityType2PropertyRepo.(entity_type2property.Repository)
		if !ok {
			return nil, errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", entity_type2property.EntityName, entity_type2property.EntityName, entityType2PropertyRepo)
		}

		r.model = relation.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewRelationRepository(r, &entityType2PropertyRepository)
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
