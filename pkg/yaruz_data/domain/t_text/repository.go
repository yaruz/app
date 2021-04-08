package t_text

import (
	"context"

	"github.com/minipkg/selection_condition"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*TText, error)
	First(ctx context.Context, entity *TText) (*TText, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]TText, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (uint, error)
	Create(ctx context.Context, entity *TText) error
	Update(ctx context.Context, entity *TText) error
	Save(ctx context.Context, entity *TText) error
	Delete(ctx context.Context, id uint) error
}
