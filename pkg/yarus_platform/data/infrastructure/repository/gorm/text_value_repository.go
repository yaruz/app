package gorm

import (
	"context"
	"errors"
	"sort"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"

	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"

	"github.com/yaruz/app/internal/pkg/apperror"

	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"
)

// TextValueRepository is a repository for the model entity
type TextValueRepository struct {
	repository
	langFinder entity.LangFinder
}

var _ text_value.Repository = (*TextValueRepository)(nil)

// New creates a new TextValueRepository
func NewTextValueRepository(repository *repository, langFinder entity.LangFinder) (*TextValueRepository, error) {
	return &TextValueRepository{
		repository: *repository,
		langFinder: langFinder,
	}, nil
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *TextValueRepository) queryTx(ctx context.Context, cond *selection_condition.SelectionCondition, tx *gorm.DB) ([]text_value.TextValue, error) {
	items := []text_value.TextValue{}
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

func (r *TextValueRepository) BatchDeleteTx(ctx context.Context, cond *selection_condition.SelectionCondition, tx *gorm.DB) (err error) {
	db := &minipkg_gorm.DB{GormDB: tx}
	if db, err = db.SchemeInitWithContext(ctx, r.model); err != nil {
		return err
	}

	gormDB := minipkg_gorm.Conditions(db.DB(), cond)
	if gormDB.Error != nil {
		return gormDB.Error
	}

	if err = gormDB.Delete(&text_value.TextValue{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}

func (r *TextValueRepository) BatchSaveChangesTx(ctx context.Context, entityID uint, values []text_value.TextValue, langID uint, tx *gorm.DB) (err error) {
	db := &minipkg_gorm.DB{GormDB: tx}
	if db, err = db.SchemeInitWithContext(ctx, r.model); err != nil {
		return err
	}

	return db.DB().Transaction(func(tx *gorm.DB) error {
		oldValues, err := r.queryTx(ctx, &selection_condition.SelectionCondition{
			Where: &text_value.TextValue{
				EntityID: entityID,
				LangID:   langID,
			},
		}, tx)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			oldValues = make([]text_value.TextValue, 0)
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

		return r.BatchSaveTx(ctx, entityID, values, langID, tx)
	})
}

func (r *TextValueRepository) BatchSaveTx(ctx context.Context, entityID uint, values []text_value.TextValue, langID uint, tx *gorm.DB) error {
	if len(values) == 0 {
		return nil
	}

	cfgname, err := r.langFinder.GetCfgnameByID(ctx, langID)
	if err != nil {
		return err
	}

	for i, value := range values {

		if value.ID == 0 {
			err = tx.Raw("INSERT INTO "+text_value.TableName+" (entity_id, lang_id, property_id, value, value_tsvector) VALUES (?, ?, ?, ?, ?) RETURNING id", entityID, langID, value.PropertyID, value.Value, gorm.Expr("to_tsvector(?, ?)", cfgname, value.Value)).Scan(&values[i].ID).Error
		} else {
			err = tx.Exec("UPDATE "+text_value.TableName+" SET entity_id = ?, lang_id = ?, property_id = ?, value = ?, value_tsvector = ? WHERE id = ?", entityID, langID, value.PropertyID, value.Value, gorm.Expr("to_tsvector(?, ?)", cfgname, value.Value), value.ID).Error
		}

		if err != nil {
			return err
		}
	}
	return nil
}
