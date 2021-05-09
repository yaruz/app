package entity_type

import (
	"context"

	"github.com/minipkg/selection_condition"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*EntityType, error)
	First(ctx context.Context, entity *EntityType) (*EntityType, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]EntityType, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *EntityType) error
	Update(ctx context.Context, entity *EntityType) error
	Save(ctx context.Context, entity *EntityType) error
	Delete(ctx context.Context, id uint) error
}
