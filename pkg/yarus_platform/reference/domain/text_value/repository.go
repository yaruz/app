package text_value

import (
	"context"

	"gorm.io/gorm"

	"github.com/minipkg/selection_condition"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*TextValue, error)
	First(ctx context.Context, entity *TextValue) (*TextValue, error)
	FirstTx(ctx context.Context, tx *gorm.DB, entity *TextValue) (*TextValue, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]TextValue, error)
	QueryTx(ctx context.Context, tx *gorm.DB, cond *selection_condition.SelectionCondition) ([]TextValue, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *TextValue) error
	CreateTx(ctx context.Context, tx *gorm.DB, entity *TextValue) error
	Update(ctx context.Context, entity *TextValue) error
	UpdateTx(ctx context.Context, tx *gorm.DB, entity *TextValue) error
	Save(ctx context.Context, entity *TextValue) error
	Delete(ctx context.Context, entity *TextValue) error
	DeleteTx(ctx context.Context, tx *gorm.DB, entity *TextValue) error
	GetValuesTx(ctx context.Context, tx *gorm.DB, langID uint, sourceIDs ...*uint) ([]*string, error)
}
