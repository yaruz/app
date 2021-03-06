package text_value

import (
	"context"

	"github.com/minipkg/selection_condition"
	"gorm.io/gorm"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	//SetDefaultConditions(conditions *selection_condition.SelectionCondition)
	// Get returns the album with the specified album ID.
	//Get(ctx context.Context, id uint) (*TextValue, error)
	//First(ctx context.Context, entity *TextValue) (*TextValue, error)
	// Query returns the list of albums with the given offset and limit.
	//Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]TextValue, error)
	//Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	//Create(ctx context.Context, entity *TextValue) error
	//Update(ctx context.Context, entity *TextValue) error
	//Save(ctx context.Context, entity *TextValue) error
	//Delete(ctx context.Context, id uint) error
	BatchSaveChangesTx(ctx context.Context, entityID uint, values []TextValue, langID uint, tx *gorm.DB) error
	BatchDeleteTx(ctx context.Context, cond *selection_condition.SelectionCondition, tx *gorm.DB) error
}
