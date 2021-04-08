package property_type

import (
	"context"

	"github.com/minipkg/selection_condition"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*PropertyType, error)
	First(ctx context.Context, entity *PropertyType) (*PropertyType, error)
	// Query returns the list of items with the given conditions.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]PropertyType, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (uint, error)
	Create(ctx context.Context, entity *PropertyType) error
	Update(ctx context.Context, entity *PropertyType) error
	Save(ctx context.Context, entity *PropertyType) error
	Delete(ctx context.Context, id uint) error
}
