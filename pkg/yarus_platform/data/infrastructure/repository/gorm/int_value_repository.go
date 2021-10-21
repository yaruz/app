package gorm

import (
	"context"
	"errors"
	"sort"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/int_value"

	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"

	"github.com/yaruz/app/internal/pkg/apperror"

	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"
)

// IntValueRepository is a repository for the model entity
type IntValueRepository struct {
	repository
}

var _ int_value.Repository = (*IntValueRepository)(nil)

// New creates a new IntValueRepository
func NewIntValueRepository(repository *repository) (*IntValueRepository, error) {
	return &IntValueRepository{repository: *repository}, nil
}

// Query retrieves the records with the specified offset and limit from the database.
func (r *IntValueRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]int_value.IntValue, error) {
	items := []int_value.IntValue{}
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

func (r *IntValueRepository) BatchDeleteTx(ctx context.Context, cond *selection_condition.SelectionCondition, tx *gorm.DB) error {
	db := minipkg_gorm.Conditions(r.db.GormTx(tx), cond)
	if db.Error != nil {
		return db.Error
	}

	err := db.Delete(&int_value.IntValue{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}

func (r *IntValueRepository) BatchSaveChangesTx(ctx context.Context, entityID uint, values []int_value.IntValue, langID uint, tx *gorm.DB) error {
	return r.db.GormTx(tx).Transaction(func(tx *gorm.DB) error {
		oldValues, err := r.Query(ctx, &selection_condition.SelectionCondition{
			Where: &int_value.IntValue{
				EntityID: entityID,
			},
		})
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			oldValues = make([]int_value.IntValue, 0)
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
