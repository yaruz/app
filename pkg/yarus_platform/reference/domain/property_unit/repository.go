package property_unit

import (
	"context"

	"github.com/minipkg/selection_condition"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*PropertyUnit, error)
	First(ctx context.Context, entity *PropertyUnit) (*PropertyUnit, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]PropertyUnit, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *PropertyUnit) error
	Update(ctx context.Context, entity *PropertyUnit) error
	Delete(ctx context.Context, entity *PropertyUnit) error
	TGet(ctx context.Context, id uint, langID uint) (*PropertyUnit, error)
	TFirst(ctx context.Context, entity *PropertyUnit, langID uint) (*PropertyUnit, error)
	TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]PropertyUnit, error)
	TCreate(ctx context.Context, entity *PropertyUnit, langID uint) (err error)
	TUpdate(ctx context.Context, entity *PropertyUnit, langID uint) (err error)
}
