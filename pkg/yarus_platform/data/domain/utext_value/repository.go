package utext_value

import (
	"context"

	"gorm.io/gorm"

	"github.com/minipkg/selection_condition"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	BatchSaveChangesTx(ctx context.Context, entityID uint, values []UTextValue, langID uint, tx *gorm.DB) error
	BatchDeleteTx(ctx context.Context, cond *selection_condition.SelectionCondition, tx *gorm.DB) error
}
