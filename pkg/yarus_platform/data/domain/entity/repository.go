package entity

import (
	"context"

	"github.com/minipkg/selection_condition"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*Entity, error)
	First(ctx context.Context, entity *Entity) (*Entity, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]Entity, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *Entity) error
	Update(ctx context.Context, entity *Entity) error
	Save(ctx context.Context, entity *Entity) error
	Delete(ctx context.Context, id uint) error
}
