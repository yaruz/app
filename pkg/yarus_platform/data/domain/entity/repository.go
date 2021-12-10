package entity

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/bool_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/date_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/float_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/int_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/time_value"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"
)

type ValueRepositories struct {
	Text  text_value.Repository
	Bool  bool_value.Repository
	Int   int_value.Repository
	Float float_value.Repository
	Date  date_value.Repository
	Time  time_value.Repository
}

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	//SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	//Get(ctx context.Context, id uint, typeID uint, langID uint) (*Entity, error)
	//First(ctx context.Context, entity *Entity, langID uint) (*Entity, error)
	// Query returns the list of albums with the given offset and limit.
	//Query(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]Entity, error)
	//Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *Entity, langID uint) error
	Update(ctx context.Context, entity *Entity, langID uint) error
	Delete(ctx context.Context, id uint, typeID uint) error
}

type EntityIDRepository interface {
	AutoMigrate() error
	NextVal(entityTypeSysname string) (id uint, err error)
	LastVal(entityTypeSysname string) (id uint, err error)
}
