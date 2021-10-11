package date_value

import (
	"context"
	"time"

	"github.com/minipkg/selection_condition"
	"gorm.io/gorm"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	BatchSaveChangesTx(ctx context.Context, entityID uint, mapOfValues map[uint]time.Time, tx *gorm.DB) error
	BatchDeleteTx(ctx context.Context, cond *selection_condition.SelectionCondition, tx *gorm.DB) error
}
