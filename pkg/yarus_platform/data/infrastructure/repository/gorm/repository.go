package gorm

import (
	"github.com/pkg/errors"
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
		textValueRepository, err := r.getTextValueRepository(logger, dbase)
		if err != nil {
			return nil, err
		}
		repo, err = NewEntityRepository(r, textValueRepository)
	case text_value.EntityName:
		repo, err = NewTextValueRepository(r)
	default:
		err = errors.Errorf("Text for entity %q not found", entityName)
	}
	return repo, err
}

func (r *repository) getTextValueRepository(logger log.ILogger, dbase minipkg_gorm.IDB) (text_value.Repository, error) {
	textValueRepo, err := GetRepository(logger, dbase, text_value.EntityName)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not get db repository for entity %q, error happened: %v", text_value.EntityName, err)
	}
	textValueRepository, ok := textValueRepo.(text_value.Repository)
	if !ok {
		return nil, errors.Errorf("Can not cast DB repository for entity %q to %vRepository. Repo: %v", text_value.EntityName, text_value.EntityName, textValueRepo)
	}

	return textValueRepository, nil
}

func (r repository) DB() *gorm.DB {
	return minipkg_gorm.Conditions(r.db.DB(), r.Conditions)
}
