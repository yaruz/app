package gorm

import (
	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_value"

	"github.com/pkg/errors"
	"gorm.io/gorm"

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
