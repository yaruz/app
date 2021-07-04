package entity_type2property

import (
	"context"

	"github.com/minipkg/selection_condition"
	"gorm.io/gorm"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*EntityType2Property, error)
	First(ctx context.Context, entity *EntityType2Property) (*EntityType2Property, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]EntityType2Property, error)
	QueryWithEntityType(ctx context.Context, cond *selection_condition.SelectionCondition) ([]EntityType2Property, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *EntityType2Property) error
	CreateTx(ctx context.Context, tx *gorm.DB, entity *EntityType2Property) error
	Delete(ctx context.Context, entity *EntityType2Property) error
	DeleteTx(ctx context.Context, tx *gorm.DB, entity *EntityType2Property) error
}
