package bool_value

import (
	"context"

	"github.com/minipkg/selection_condition"
	"gorm.io/gorm"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	BatchSaveChangesTx(ctx context.Context, entityID uint, values []BoolValue, langID uint, tx *gorm.DB) error
	BatchDeleteTx(ctx context.Context, cond *selection_condition.SelectionCondition, tx *gorm.DB) error
}
