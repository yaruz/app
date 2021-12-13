package gorm

import (
	"context"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"
)

type IMapReducer interface {
	GetDB(ctx context.Context, typeID uint, ID uint) (minipkg_gorm.IDB, error)
	GetDBs(entityWhereConditions selection_condition.WhereConditions) []minipkg_gorm.IDB
	Query(ctx context.Context, model interface{}, entityWhereConditions selection_condition.WhereConditions, f func(db minipkg_gorm.IDB) ([]SearchResult, error)) ([]SearchResult, error)
	Count(ctx context.Context, model interface{}, entityWhereConditions selection_condition.WhereConditions, f func(db minipkg_gorm.IDB) (uint, error)) (uint, error)
}
