package gorm

import (
	"context"
	"errors"
	"sort"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/time_value"

	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"

	"github.com/yaruz/app/internal/pkg/apperror"

	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"
)

// TimeValueRepository is a repository for the model entity
type TimeValueRepository struct {
	repository
}

var _ time_value.Repository = (*TimeValueRepository)(nil)

// New creates a new TimeValueRepository
func NewTimeValueRepository(repository *repository) (*TimeValueRepository, error) {
	return &TimeValueRepository{repository: *repository}, nil
}

// Query retrieves the records with the specified offset and limit from the database.
func (r *TimeValueRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]time_value.TimeValue, error) {
	items := []time_value.TimeValue{}
	db := minipkg_gorm.Conditions(r.DB(), cond)
	if db.Error != nil {
		return nil, db.Error
	}

	err := db.Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return items, yaruserror.ErrNotFound
		}
		return nil, err
	}
	return items, err
}

func (r *TimeValueRepository) BatchDeleteTx(ctx context.Context, cond *selection_condition.SelectionCondition, tx *gorm.DB) error {
	db := minipkg_gorm.Conditions(r.db.GormTx(tx), cond)
	if db.Error != nil {
		return db.Error
	}

	err := db.Delete(&time_value.TimeValue{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}

func (r *TimeValueRepository) BatchSaveChangesTx(ctx context.Context, entityID uint, values []time_value.TimeValue, langID uint, tx *gorm.DB) error {
	for i := range values {
		values[i].EntityID = entityID
	}
	return r.db.GormTx(tx).Transaction(func(tx *gorm.DB) error {
		oldValues, err := r.Query(ctx, &selection_condition.SelectionCondition{
			Where: &time_value.TimeValue{
				EntityID: entityID,
			},
		})
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			oldValues = make([]time_value.TimeValue, 0)
		}

		oldValuesIds := make([]int, 0, len(oldValues))
		for _, oldValue := range oldValues {
			oldValuesIds = append(oldValuesIds, int(oldValue.ID))
		}
		sort.Ints(oldValuesIds)

		for _, value := range values {
			i := sort.SearchInts(oldValuesIds, int(value.ID))
			if i < len(oldValuesIds) && oldValuesIds[i] == int(value.ID) {
				oldValuesIds = append(oldValuesIds[:i], oldValuesIds[i+1:]...)
			}
		}

		if len(oldValuesIds) > 0 {
			if err := r.BatchDeleteTx(ctx, &selection_condition.SelectionCondition{
				Where: selection_condition.WhereCondition{
					Field:     "ID",
					Condition: selection_condition.ConditionIn,
					Value:     selection_condition.IntSlice2EmptyInterfaceSlice(oldValuesIds),
				},
			}, tx); err != nil {
				return err
			}
		}

		if len(values) > 0 {
			return tx.Save(values).Error
		}

		return nil
	})
}
