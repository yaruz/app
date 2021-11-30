package gorm

import (
	"context"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"
)

type IMapReducer interface {
	GetDB(ID uint, typeID uint) minipkg_gorm.IDB
	GetDBs(condition *selection_condition.SelectionCondition) []minipkg_gorm.IDB
	GetDBForInsert(typeID uint) minipkg_gorm.IDB
	Query(ctx context.Context, model interface{}, condition *selection_condition.SelectionCondition, f func(db minipkg_gorm.IDB) ([]SearchResult, error)) ([]SearchResult, error)
	Count(ctx context.Context, model interface{}, condition *selection_condition.SelectionCondition, f func(db minipkg_gorm.IDB) (uint, error)) (uint, error)
}
