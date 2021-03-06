package property

import (
	"context"

	"github.com/minipkg/selection_condition"
	"gorm.io/gorm"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*Property, error)
	GetBySysname(ctx context.Context, sysname string, langID uint) (*Property, error)
	First(ctx context.Context, entity *Property) (*Property, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]Property, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *Property) error
	Update(ctx context.Context, entity *Property) error
	Delete(ctx context.Context, entity *Property) error
	TGet(ctx context.Context, id uint, langID uint) (*Property, error)
	TFirst(ctx context.Context, entity *Property, langID uint) (*Property, error)
	TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]Property, error)
	TCreate(ctx context.Context, entity *Property, langID uint) (err error)
	TUpdate(ctx context.Context, entity *Property, langID uint) (err error)
	EntityNameAndDescriptionInitTx(ctx context.Context, tx *gorm.DB, entity *Property, langID uint) error
}
