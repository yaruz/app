package gorm

import (
	"context"
	"errors"

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
	r := &BoolValueRepository{repository: *repository}
	r.autoMigrate()
	return r, nil
}

func (r *BoolValueRepository) autoMigrate() {
	if r.db.IsAutoMigrate() {
		r.db.DB().AutoMigrate(&bool_value.BoolValue{})
	}
}

// Query retrieves the album records with the specified offset and limit from the database.
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
	db := minipkg_gorm.Conditions(tx, cond)
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

func (r *BoolValueRepository) BatchSaveChangesTx(ctx context.Context, entityID uint, mapOfValues map[uint]bool, tx *gorm.DB) error {
	return tx.Transaction(func(tx *gorm.DB) error {
		var valueObj *bool_value.BoolValue
		// можно и без этого запроса, а просто брать из entity.BoolValues, но для большей безопасности сделаем отдельный независимый запрос
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

		mapOldValues := make(map[uint]*bool_value.BoolValue, len(oldValues))
		for i := range oldValues {
			mapOldValues[oldValues[i].PropertyID] = &oldValues[i]
		}

		newValues := make([]bool_value.BoolValue, 0, len(oldValues))
		for propertyID, value := range mapOfValues {
			if _, ok := mapOldValues[propertyID]; ok {
				valueObj = mapOldValues[propertyID]
				delete(mapOldValues, propertyID)
			} else {
				valueObj = &bool_value.BoolValue{
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
