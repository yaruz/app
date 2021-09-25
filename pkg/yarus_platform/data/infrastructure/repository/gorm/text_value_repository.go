package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"

	"github.com/yaruz/app/internal/pkg/apperror"

	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// TextValueRepository is a repository for the model entity
type TextValueRepository struct {
	repository
}

var _ text_value.Repository = (*TextValueRepository)(nil)

// New creates a new TextValueRepository
func NewTextValueRepository(repository *repository) (*TextValueRepository, error) {
	r := &TextValueRepository{repository: *repository}
	r.autoMigrate()
	return r, nil
}

func (r *TextValueRepository) autoMigrate() {
	if r.db.IsAutoMigrate() {
		r.db.DB().AutoMigrate(&text_value.TextValue{})
	}
}

// Get reads the album with the specified ID from the database.
func (r *TextValueRepository) Get(ctx context.Context, id uint) (*text_value.TextValue, error) {
	entity := &text_value.TextValue{}

	err := r.DB().First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

func (r *TextValueRepository) First(ctx context.Context, entity *text_value.TextValue) (*text_value.TextValue, error) {
	err := r.DB().Where(entity).First(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *TextValueRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]text_value.TextValue, error) {
	items := []text_value.TextValue{}
	db := minipkg_gorm.Conditions(r.DB(), cond)
	if db.Error != nil {
		return nil, db.Error
	}

	err := db.Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return items, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return items, err
}

func (r *TextValueRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
	var count int64
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.DB(), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	err := db.Count(&count).Error
	return count, err
}

// Create saves a new record in the database.
func (r *TextValueRepository) Create(ctx context.Context, entity *text_value.TextValue) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}

	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *TextValueRepository) Update(ctx context.Context, entity *text_value.TextValue) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}

	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *TextValueRepository) Save(ctx context.Context, entity *text_value.TextValue) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *TextValueRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(&text_value.TextValue{}, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}

func (r *TextValueRepository) BatchDeleteTx(ctx context.Context, cond *selection_condition.SelectionCondition, tx *gorm.DB) error {
	db := minipkg_gorm.Conditions(tx, cond)
	if db.Error != nil {
		return db.Error
	}

	err := db.Delete(&text_value.TextValue{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}

func (r *TextValueRepository) BatchSaveChangesTx(ctx context.Context, entityID uint, langID uint, mapOfStrings map[uint]string, tx *gorm.DB) error {
	return tx.Transaction(func(tx *gorm.DB) error {
		var textValue *text_value.TextValue
		// можно и без этого запроса, а просто брать из entity.TextValues, но для большей безопасности сделаем отдельный независимый запрос
		oldTextValues, err := r.Query(ctx, &selection_condition.SelectionCondition{
			Where: &text_value.TextValue{
				EntityID: entityID,
				LangID:   langID,
			},
		})
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			oldTextValues = make([]text_value.TextValue, 0)
		}

		mapOldTextValues := make(map[uint]*text_value.TextValue, len(oldTextValues))
		for i := range oldTextValues {
			mapOldTextValues[oldTextValues[i].PropertyID] = &oldTextValues[i]
		}

		newTextValues := make([]text_value.TextValue, 0, len(oldTextValues))
		for propertyID, value := range mapOfStrings {
			if _, ok := mapOldTextValues[propertyID]; ok {
				textValue = mapOldTextValues[propertyID]
				delete(mapOldTextValues, propertyID)
			} else {
				textValue = &text_value.TextValue{
					EntityID:   entityID,
					LangID:     langID,
					PropertyID: propertyID,
				}
			}
			textValue.Value = value
			newTextValues = append(newTextValues, *textValue)
		}

		if len(mapOldTextValues) > 0 {
			delTextValuesIds := make([]uint, 0)
			for _, t := range mapOldTextValues {
				delTextValuesIds = append(delTextValuesIds, t.ID)
			}

			if err := r.BatchDeleteTx(ctx, &selection_condition.SelectionCondition{
				Where: selection_condition.WhereCondition{
					Field:     "id",
					Condition: selection_condition.ConditionIn,
					Value:     delTextValuesIds,
				},
			}, tx); err != nil {
				return err
			}
		}

		return tx.Save(newTextValues).Error
	})
}
