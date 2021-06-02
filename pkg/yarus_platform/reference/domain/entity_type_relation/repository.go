package entity_type_relation

import (
	"context"

	"github.com/minipkg/selection_condition"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*EntityTypeRelation, error)
	First(ctx context.Context, entity *EntityTypeRelation) (*EntityTypeRelation, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]EntityTypeRelation, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *EntityTypeRelation) error
	Update(ctx context.Context, entity *EntityTypeRelation) error
	Save(ctx context.Context, entity *EntityTypeRelation) error
	Delete(ctx context.Context, id uint) error
}
