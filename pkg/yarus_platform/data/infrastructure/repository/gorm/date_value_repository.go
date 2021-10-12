package gorm

import (
	"context"
	"errors"
	"time"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/date_value"

	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"

	"github.com/yaruz/app/internal/pkg/apperror"

	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"
)

// DateValueRepository is a repository for the model entity
type DateValueRepository struct {
	repository
}

var _ date_value.Repository = (*DateValueRepository)(nil)

// New creates a new DateValueRepository
func NewDateValueRepository(repository *repository) (*DateValueRepository, error) {
	return &DateValueRepository{repository: *repository}, nil
}

// Query retrieves the records with the specified offset and limit from the database.
func (r *DateValueRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]date_value.DateValue, error) {
	items := []date_value.DateValue{}
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

func (r *DateValueRepository) BatchDeleteTx(ctx context.Context, cond *selection_condition.SelectionCondition, tx *gorm.DB) error {
	db := minipkg_gorm.Conditions(tx, cond)
	if db.Error != nil {
		return db.Error
	}

	err := db.Delete(&date_value.DateValue{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}

func (r *DateValueRepository) BatchSaveChangesTx(ctx context.Context, entityID uint, mapOfValues map[uint]time.Time, tx *gorm.DB) error {
	return tx.Transaction(func(tx *gorm.DB) error {
		var valueObj *date_value.DateValue
		// можно и без этого запроса, а просто брать из entity.DateValues, но для большей безопасности сделаем отдельный независимый запрос
		oldValues, err := r.Query(ctx, &selection_condition.SelectionCondition{
			Where: &date_value.DateValue{
				EntityID: entityID,
			},
		})
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			oldValues = make([]date_value.DateValue, 0)
		}

		mapOldValues := make(map[uint]*date_value.DateValue, len(oldValues))
		for i := range oldValues {
			mapOldValues[oldValues[i].PropertyID] = &oldValues[i]
		}

		newValues := make([]date_value.DateValue, 0, len(oldValues))
		for propertyID, value := range mapOfValues {
			if _, ok := mapOldValues[propertyID]; ok {
				valueObj = mapOldValues[propertyID]
				delete(mapOldValues, propertyID)
			} else {
				valueObj = &date_value.DateValue{
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
