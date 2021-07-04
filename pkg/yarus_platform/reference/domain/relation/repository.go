package relation

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"

	"github.com/minipkg/selection_condition"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*Relation, error)
	First(ctx context.Context, entity *Relation) (*Relation, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]Relation, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *Relation) error
	Update(ctx context.Context, entity *Relation) error
	Save(ctx context.Context, entity *Relation) error
	Delete(ctx context.Context, id uint) error
	PropertyAndRelationQuery(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property.Property, []Relation, error)
	GetPropertiesAndRelationsByEntityTypeID(ctx context.Context, entityTypeID uint) ([]property.Property, []Relation, error)
}
