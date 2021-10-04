package gorm

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_lang"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/log"
	"github.com/minipkg/selection_condition"
	"github.com/pkg/errors"
	"gorm.io/gorm"

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
	var textSourceRepository text_source.Repository

	switch entityName {
	case entity_type.EntityName, property.EntityName, entity_type.RelationEntityName, property_type.EntityName, property_view_type.EntityName, property_unit.EntityName, property_group.EntityName:
		if textSourceRepository, err = r.getTextSourceRepository(logger, dbase); err != nil {
			return nil, err
		}
	}

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
		repo, err = NewEntityTypeRepository(r, &entityType2PropertyRepository, textSourceRepository)
	case property_group.EntityName:
		r.model = property_group.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewPropertyGroupRepository(r, textSourceRepository)
	case property.EntityName:
		r.model = property.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewPropertyRepository(r, textSourceRepository)
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
		repo, err = NewPropertyTypeRepository(r, textSourceRepository)
	case property_unit.EntityName:
		r.model = property_unit.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewPropertyUnitRepository(r, textSourceRepository)
	case property_view_type.EntityName:
		r.model = property_view_type.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewPropertyViewTypeRepository(r, textSourceRepository)
	case text_source.EntityName:
		textValueRepo, err := GetRepository(logger, dbase, text_value.EntityName)
		if err != nil {
			return nil, errors.Wrapf(err, "Can not get db repository for entity %q, error happened: %v", text_value.EntityName, err)
		}
		textValueRepository, ok := textValueRepo.(text_value.Repository)
		if !ok {
			return nil, errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", text_value.EntityName, text_value.EntityName, textValueRepo)
		}

		r.model = text_source.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewTextSourceRepository(r, textValueRepository)
	case text_value.EntityName:
		r.model = text_value.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewTextValueRepository(r)
	case text_lang.EntityName:
		r.model = text_lang.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewTextLangRepository(r)
	case entity_type.RelationEntityName:
		propertyRepo, err := GetRepository(logger, dbase, property.EntityName)
		if err != nil {
			return nil, errors.Wrapf(err, "Can not get db repository for entity %q, error happened: %v", property.EntityName, err)
		}
		propertyRepository, ok := propertyRepo.(property.Repository)
		if !ok {
			return nil, errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", property.EntityName, property.EntityName, propertyRepo)
		}

		entityType2PropertyRepo, err := GetRepository(logger, dbase, entity_type2property.EntityName)
		if err != nil {
			return nil, errors.Wrapf(err, "Can not get db repository for entity %q, error happened: %v", entity_type2property.EntityName, err)
		}
		entityType2PropertyRepository, ok := entityType2PropertyRepo.(entity_type2property.Repository)
		if !ok {
			return nil, errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", entity_type2property.EntityName, entity_type2property.EntityName, entityType2PropertyRepo)
		}

		r.model = entity_type.NewRelation()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewRelationRepository(r, propertyRepository, entityType2PropertyRepository, textSourceRepository)
	default:
		err = errors.Errorf("Repository for entity %q not found", entityName)
	}
	return repo, err
}

func (r *repository) getTextSourceRepository(logger log.ILogger, dbase minipkg_gorm.IDB) (text_source.Repository, error) {
	textSourceRepo, err := GetRepository(logger, dbase, text_source.EntityName)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not get db repository for entity %q, error happened: %v", text_source.EntityName, err)
	}
	textSourceRepository, ok := textSourceRepo.(text_source.Repository)
	if !ok {
		return nil, errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", text_source.EntityName, text_source.EntityName, textSourceRepo)
	}

	return textSourceRepository, nil
}

func (r *repository) SetDefaultConditions(defaultConditions *selection_condition.SelectionCondition) {
	r.Conditions = defaultConditions
}

func (r repository) DB() *gorm.DB {
	return minipkg_gorm.Conditions(r.db.DB(), r.Conditions)
}
