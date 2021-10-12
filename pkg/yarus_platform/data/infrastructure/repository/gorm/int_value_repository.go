package gorm

import (
	"context"
	"errors"

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
	db := minipkg_gorm.Conditions(tx, cond)
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

func (r *IntValueRepository) BatchSaveChangesTx(ctx context.Context, entityID uint, mapOfValues map[uint]int, tx *gorm.DB) error {
	return tx.Transaction(func(tx *gorm.DB) error {
		var valueObj *int_value.IntValue
		// можно и без этого запроса, а просто брать из entity.IntValues, но для большей безопасности сделаем отдельный независимый запрос
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

		mapOldValues := make(map[uint]*int_value.IntValue, len(oldValues))
		for i := range oldValues {
			mapOldValues[oldValues[i].PropertyID] = &oldValues[i]
		}

		newValues := make([]int_value.IntValue, 0, len(oldValues))
		for propertyID, value := range mapOfValues {
			if _, ok := mapOldValues[propertyID]; ok {
				valueObj = mapOldValues[propertyID]
				delete(mapOldValues, propertyID)
			} else {
				valueObj = &int_value.IntValue{
					EntityID:   entityID,
					PropertyID: propertyID,
				}
			}
			valueObj.Value = value
			newValues = append(newValues, *valueObj)
		}

		if len(mapOldValues) > 0 {
			delValuesIds := make([]uint, 0)
			for _, t := range mapOldValues {
				delValuesIds = append(delValuesIds, t.ID)
			}

			if err := r.BatchDeleteTx(ctx, &selection_condition.SelectionCondition{
				Where: selection_condition.WhereCondition{
					Field:     "id",
					Condition: selection_condition.ConditionIn,
					Value:     delValuesIds,
				},
			}, tx); err != nil {
				return err
			}
		}

		if len(newValues) > 0 {
			return tx.Save(newValues).Error
		}
		return nil
	})
}
