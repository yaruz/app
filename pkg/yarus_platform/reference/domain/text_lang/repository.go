package text_lang

import (
	"context"

	"github.com/minipkg/selection_condition"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, id uint) (*TextLang, error)
	First(ctx context.Context, entity *TextLang) (*TextLang, error)
	// Query returns the list of albums with the given offset and limit.
	Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]TextLang, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *TextLang) error
	Upsert(ctx context.Context, entity *TextLang) error
	Update(ctx context.Context, entity *TextLang) error
	Save(ctx context.Context, entity *TextLang) error
	Delete(ctx context.Context, id uint) error
}
