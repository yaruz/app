package gorm

import (
	"context"
	"errors"
	"sort"

	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/utext_value"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
)

// UTextValueRepository is a repository for the model entity
type UTextValueRepository struct {
	repository
}

var _ utext_value.Repository = (*UTextValueRepository)(nil)

// NewUTextValueRepository creates a new UTextValueRepository
func NewUTextValueRepository(repository *repository) (*UTextValueRepository, error) {
	return &UTextValueRepository{repository: *repository}, nil
}

// Query retrieves the records with the specified offset and limit from the database.
func (r *UTextValueRepository) queryTx(ctx context.Context, cond *selection_condition.SelectionCondition, tx *gorm.DB) ([]utext_value.UTextValue, error) {
	items := []utext_value.UTextValue{}
	db := minipkg_gorm.Conditions(tx, cond)
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

func (r *UTextValueRepository) BatchDeleteTx(ctx context.Context, cond *selection_condition.SelectionCondition, tx *gorm.DB) (err error) {
	db := &minipkg_gorm.DB{GormDB: tx}
	if db, err = db.SchemeInitWithContext(ctx, r.model); err != nil {
		return err
	}

	gormDB := minipkg_gorm.Conditions(db.DB(), cond)
	if gormDB.Error != nil {
		return gormDB.Error
	}

	if err = gormDB.Delete(&utext_value.UTextValue{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}

func (r *UTextValueRepository) BatchSaveChangesTx(ctx context.Context, entityID uint, values []utext_value.UTextValue, langID uint, tx *gorm.DB) (err error) {
	for i := range values {
		values[i].EntityID = entityID
	}

	db := &minipkg_gorm.DB{GormDB: tx}
	if db, err = db.SchemeInitWithContext(ctx, r.model); err != nil {
		return err
	}

	return db.DB().Transaction(func(tx *gorm.DB) error {
		oldValues, err := r.queryTx(ctx, &selection_condition.SelectionCondition{
			Where: &utext_value.UTextValue{
				EntityID: entityID,
			},
		}, tx)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			oldValues = make([]utext_value.UTextValue, 0)
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
