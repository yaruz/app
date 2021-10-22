package gorm

import (
	"context"
	"errors"
	"sort"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/bool_value"

	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"

	"github.com/yaruz/app/internal/pkg/apperror"

	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"
)

// BoolValueRepository is a repository for the model entity
type BoolValueRepository struct {
	repository
}

var _ bool_value.Repository = (*BoolValueRepository)(nil)

// New creates a new BoolValueRepository
func NewBoolValueRepository(repository *repository) (*BoolValueRepository, error) {
	return &BoolValueRepository{repository: *repository}, nil
}

// Query retrieves the records with the specified offset and limit from the database.
func (r *BoolValueRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]bool_value.BoolValue, error) {
	items := []bool_value.BoolValue{}
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

func (r *BoolValueRepository) BatchDeleteTx(ctx context.Context, cond *selection_condition.SelectionCondition, tx *gorm.DB) error {
	db := minipkg_gorm.Conditions(r.db.GormTx(tx), cond)
	if db.Error != nil {
		return db.Error
	}

	err := db.Delete(&bool_value.BoolValue{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}

func (r *BoolValueRepository) BatchSaveChangesTx(ctx context.Context, entityID uint, values []bool_value.BoolValue, langID uint, tx *gorm.DB) error {
	for i := range values {
		values[i].EntityID = entityID
	}
	return r.db.GormTx(tx).Transaction(func(tx *gorm.DB) error {
		oldValues, err := r.Query(ctx, &selection_condition.SelectionCondition{
			Where: &bool_value.BoolValue{
				EntityID: entityID,
			},
		})
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			oldValues = make([]bool_value.BoolValue, 0)
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
