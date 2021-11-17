package gorm

import (
	"context"

	"github.com/pkg/errors"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/bool_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/date_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/float_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/int_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/time_value"
	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/log"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"
)

// IRepository is an interface of repository
type IRepository interface{}

// repository persists albums in database
type repository struct {
	db         minipkg_gorm.IDB
	logger     log.ILogger
	Conditions *selection_condition.SelectionCondition
	model      interface{}
}

const DefaultLimit = 1000

// GetRepository return a repository
func GetRepository(logger log.ILogger, dbase minipkg_gorm.IDB, entityName string, langFinder entity.LangFinder) (repo IRepository, err error) {
	r := &repository{
		logger: logger,
	}
	ctx := context.Background()

	switch entityName {
	case entity.EntityName:
		valueRepositories, err := r.getValueRepositories(logger, dbase, langFinder)
		if err != nil {
			return nil, err
		}

		r.model = entity.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewEntityRepository(r, valueRepositories)
	case text_value.EntityName:

		r.model = text_value.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewTextValueRepository(r, langFinder)
	case bool_value.EntityName:

		r.model = bool_value.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewBoolValueRepository(r)
	case int_value.EntityName:

		r.model = int_value.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewIntValueRepository(r)
	case float_value.EntityName:

		r.model = float_value.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewFloatValueRepository(r)
	case date_value.EntityName:

		r.model = date_value.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewDateValueRepository(r)
	case time_value.EntityName:

		r.model = time_value.New()

		if r.db, err = dbase.SchemeInitWithContext(ctx, r.model); err != nil {
			return nil, err
		}
		repo, err = NewTimeValueRepository(r)
	default:
		err = errors.Errorf("Text for entity %q not found", entityName)
	}
	return repo, err
}

func (r *repository) getValueRepositories(logger log.ILogger, dbase minipkg_gorm.IDB, langFinder entity.LangFinder) (*entity.ValueRepositories, error) {
	textValueRepo, err := GetRepository(logger, dbase, text_value.EntityName, langFinder)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not get db repository for entity %q, error happened: %v", text_value.EntityName, err)
	}
	textValueRepository, ok := textValueRepo.(text_value.Repository)
	if !ok {
		return nil, errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", text_value.EntityName, text_value.EntityName, textValueRepo)
	}

	boolValueRepo, err := GetRepository(logger, dbase, bool_value.EntityName, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not get db repository for entity %q, error happened: %v", bool_value.EntityName, err)
	}
	boolValueRepository, ok := boolValueRepo.(bool_value.Repository)
	if !ok {
		return nil, errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", bool_value.EntityName, bool_value.EntityName, boolValueRepo)
	}

	intValueRepo, err := GetRepository(logger, dbase, int_value.EntityName, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not get db repository for entity %q, error happened: %v", int_value.EntityName, err)
	}
	intValueRepository, ok := intValueRepo.(int_value.Repository)
	if !ok {
		return nil, errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", int_value.EntityName, int_value.EntityName, intValueRepo)
	}

	floatValueRepo, err := GetRepository(logger, dbase, float_value.EntityName, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not get db repository for entity %q, error happened: %v", float_value.EntityName, err)
	}
	floatValueRepository, ok := floatValueRepo.(float_value.Repository)
	if !ok {
		return nil, errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", float_value.EntityName, float_value.EntityName, floatValueRepo)
	}

	dateValueRepo, err := GetRepository(logger, dbase, date_value.EntityName, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not get db repository for entity %q, error happened: %v", date_value.EntityName, err)
	}
	dateValueRepository, ok := dateValueRepo.(date_value.Repository)
	if !ok {
		return nil, errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", date_value.EntityName, date_value.EntityName, dateValueRepo)
	}

	timeValueRepo, err := GetRepository(logger, dbase, time_value.EntityName, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not get db repository for entity %q, error happened: %v", time_value.EntityName, err)
	}
	timeValueRepository, ok := timeValueRepo.(time_value.Repository)
	if !ok {
		return nil, errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", time_value.EntityName, time_value.EntityName, timeValueRepo)
	}

	return &entity.ValueRepositories{
		Text:  textValueRepository,
		Bool:  boolValueRepository,
		Int:   intValueRepository,
		Float: floatValueRepository,
		Date:  dateValueRepository,
		Time:  timeValueRepository,
	}, nil
}

func (r repository) DB() *gorm.DB {
	return minipkg_gorm.Conditions(r.db.DB(), r.Conditions)
}
