package property_group

import (
	"context"

	"github.com/minipkg/selection_condition"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*PropertyGroup, error)
	First(ctx context.Context, entity *PropertyGroup) (*PropertyGroup, error)
	// Query returns the list of items with the given conditions.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]PropertyGroup, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *PropertyGroup) error
	Update(ctx context.Context, entity *PropertyGroup) error
	Save(ctx context.Context, entity *PropertyGroup) error
	Delete(ctx context.Context, id uint) error
}
