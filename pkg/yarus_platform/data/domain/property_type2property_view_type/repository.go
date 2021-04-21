package property_type2property_view_type

import (
	"context"

	"github.com/minipkg/selection_condition"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*PropertyType2PropertyViewType, error)
	First(ctx context.Context, entity *PropertyType2PropertyViewType) (*PropertyType2PropertyViewType, error)
	// Query returns the list of items with the given conditions.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]PropertyType2PropertyViewType, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (uint, error)
	Create(ctx context.Context, entity *PropertyType2PropertyViewType) error
	Update(ctx context.Context, entity *PropertyType2PropertyViewType) error
	Save(ctx context.Context, entity *PropertyType2PropertyViewType) error
	Delete(ctx context.Context, propertyTypeID uint, propertyViewTypeID uint) error
}
