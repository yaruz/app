package property_view_type

import (
	"context"

	"github.com/minipkg/selection_condition"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*PropertyViewType, error)
	First(ctx context.Context, entity *PropertyViewType) (*PropertyViewType, error)
	// Query returns the list of items with the given conditions.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]PropertyViewType, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (uint, error)
	Create(ctx context.Context, entity *PropertyViewType) error
	Update(ctx context.Context, entity *PropertyViewType) error
	Save(ctx context.Context, entity *PropertyViewType) error
	Delete(ctx context.Context, id uint) error
}
