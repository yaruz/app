package text_source

import (
	"context"

	"github.com/minipkg/selection_condition"
	"gorm.io/gorm"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*TextSource, error)
	TGet(ctx context.Context, id uint, langID uint) (*TextSource, error)
	First(ctx context.Context, entity *TextSource) (*TextSource, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]TextSource, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *TextSource) error
	Update(ctx context.Context, entity *TextSource) error
	Save(ctx context.Context, entity *TextSource) error
	Delete(ctx context.Context, id uint) error
	DeleteTx(ctx context.Context, tx *gorm.DB, id uint) error
	GetValuesTx(ctx context.Context, tx *gorm.DB, langID uint, sourceIDs ...*uint) ([]*string, error)
	CreateValueTx(ctx context.Context, tx *gorm.DB, value *string, langID uint) (sourceID *uint, err error)
	UpdateValueTx(ctx context.Context, tx *gorm.DB, sourceID *uint, value *string, langID uint) (resSourceID *uint, err error)
}
