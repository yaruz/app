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
	GetBySysname(ctx context.Context, sysname string, langID uint) (*EntityType, error)
	First(ctx context.Context, entity *EntityType) (*EntityType, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]EntityType, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *EntityType) error
	Update(ctx context.Context, entity *EntityType) error
	Delete(ctx context.Context, entity *EntityType) error
	TGet(ctx context.Context, id uint, langID uint) (*EntityType, error)
	TFirst(ctx context.Context, entity *EntityType, langID uint) (*EntityType, error)
	TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]EntityType, error)
	TCreate(ctx context.Context, entity *EntityType, langID uint) (err error)
	TUpdate(ctx context.Context, entity *EntityType, langID uint) (err error)
	BindProperty(ctx context.Context, id uint, propertyID uint) error
	BindPropertyIfNotBinded(ctx context.Context, id uint, propertyID uint) error
	UnbindProperty(ctx context.Context, id uint, propertyID uint) error
	UnbindAllProperty(ctx context.Context, id uint) error
}
