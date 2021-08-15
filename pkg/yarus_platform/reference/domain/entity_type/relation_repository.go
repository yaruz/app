package entity_type

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"

	"github.com/minipkg/selection_condition"
)

// RelationRepository encapsulates the logic to access albums from the data source.
type RelationRepository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*Relation, error)
	First(ctx context.Context, entity *Relation) (*Relation, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]Relation, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *Relation) error
	Update(ctx context.Context, entity *Relation) error
	Delete(ctx context.Context, entity *Relation) error
	TGet(ctx context.Context, id uint, langID uint) (*Relation, error)
	TFirst(ctx context.Context, entity *Relation, langID uint) (*Relation, error)
	TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]Relation, error)
	TCreate(ctx context.Context, entity *Relation, langID uint) (err error)
	TUpdate(ctx context.Context, entity *Relation, langID uint) (err error)
	PropertyAndRelationQuery(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property.Property, []Relation, error)
	TPropertyAndRelationQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]property.Property, []Relation, error)
	GetPropertiesAndRelationsByEntityTypeID(ctx context.Context, entityTypeID uint) ([]property.Property, []Relation, error)
	TGetPropertiesAndRelationsByEntityTypeID(ctx context.Context, entityTypeID uint, langID uint) ([]property.Property, []Relation, error)
}
