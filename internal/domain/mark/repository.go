package mark

import (
	"context"
	"github.com/minipkg/selection_condition"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*Mark, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]Mark, error)
	First(ctx context.Context, entity *Mark) (*Mark, error)
}
