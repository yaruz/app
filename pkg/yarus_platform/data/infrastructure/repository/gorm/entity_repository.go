package gorm

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/bool_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/date_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/float_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/int_value"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/time_value"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/yaruz/app/internal/pkg/apperror"
	domain_entity "github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
)

// EntityRepository is a repository for the model entity
type EntityRepository struct {
	repository
	valueRepositories *domain_entity.ValueRepositories
}

var _ domain_entity.Repository = (*EntityRepository)(nil)

// New creates a new EntityRepository
func NewEntityRepository(repository *repository, valueRepositories *domain_entity.ValueRepositories) (*EntityRepository, error) {
	r := &EntityRepository{
		repository:        *repository,
		valueRepositories: valueRepositories,
	}
	r.autoMigrate()
	return r, nil
}

func (r *EntityRepository) autoMigrate() {
	if r.db.IsAutoMigrate() {
		r.db.DB().AutoMigrate(&domain_entity.Entity{})
	}
}

// Get reads the album with the specified ID from the database.
func (r *EntityRepository) Get(ctx context.Context, id uint, langID uint) (*domain_entity.Entity, error) {
	entity := &domain_entity.Entity{}

	err := r.preload(r.db.DB(), langID).First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruserror.ErrNotFound
		}
		return nil, err
	}

	return entity, err
}

func (r *EntityRepository) First(ctx context.Context, entity *domain_entity.Entity, langID uint) (*domain_entity.Entity, error) {
	err := r.preload(r.db.DB(), langID).Where(entity).First(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruserror.ErrNotFound
		}
		return nil, err
	}

	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *EntityRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]domain_entity.Entity, error) {
	items := []domain_entity.Entity{}
	db := minipkg_gorm.Conditions(r.preload(r.db.DB(), langID), cond)
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

func (r *EntityRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
	var count int64
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.db.DB(), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	err := db.Count(&count).Error
	return count, err
}

// Create saves a new record in the database.
func (r *EntityRepository) Create(ctx context.Context, entity *domain_entity.Entity, langID uint) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}

	if err := entity.Validate(); err != nil {
		return err
	}

	return r.db.DB().Transaction(func(tx *gorm.DB) (err error) {

		if err := tx.Create(entity).Error; err != nil {
			return err
		}

		return r.afterSaveTx(ctx, entity, langID, tx)
	})
}

// Update saves a changed Maintenance record in the database.
func (r *EntityRepository) Update(ctx context.Context, entity *domain_entity.Entity, langID uint) error {
	return r.updateTx(ctx, entity, langID, r.db.DB())
}

func (r *EntityRepository) updateTx(ctx context.Context, entity *domain_entity.Entity, langID uint, tx *gorm.DB) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}

	if err := entity.Validate(); err != nil {
		return err
	}

	return r.db.GormTx(tx).Transaction(func(tx *gorm.DB) (err error) {

		if err := tx.Save(entity).Error; err != nil {
			return err
		}

		return r.afterSaveTx(ctx, entity, langID, tx)
	})
}

// Сохраняем все значения всех свойств из всех слайсов ...Values
func (r *EntityRepository) afterSaveTx(ctx context.Context, entity *domain_entity.Entity, langID uint, tx *gorm.DB) error {

	if err := r.valueRepositories.Bool.BatchSaveChangesTx(ctx, entity.ID, entity.BoolValues, langID, tx); err != nil {
		return err
	}

	if err := r.valueRepositories.Int.BatchSaveChangesTx(ctx, entity.ID, entity.IntValues, langID, tx); err != nil {
		return err
	}

	if err := r.valueRepositories.Float.BatchSaveChangesTx(ctx, entity.ID, entity.FloatValues, langID, tx); err != nil {
		return err
	}

	if err := r.valueRepositories.Date.BatchSaveChangesTx(ctx, entity.ID, entity.DateValues, langID, tx); err != nil {
		return err
	}

	if err := r.valueRepositories.Time.BatchSaveChangesTx(ctx, entity.ID, entity.TimeValues, langID, tx); err != nil {
		return err
	}

	if err := r.valueRepositories.Text.BatchSaveChangesTx(ctx, entity.ID, entity.TextValues, langID, tx); err != nil {
		return err
	}

	return nil
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *EntityRepository) Delete(ctx context.Context, id uint) error {

	return r.db.DB().Transaction(func(tx *gorm.DB) (err error) {

		if err = r.valueRepositories.Bool.BatchDeleteTx(ctx, &selection_condition.SelectionCondition{
			Where: &bool_value.BoolValue{
				EntityID: id,
			},
		}, tx); err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		if err = r.valueRepositories.Int.BatchDeleteTx(ctx, &selection_condition.SelectionCondition{
			Where: &int_value.IntValue{
				EntityID: id,
			},
		}, tx); err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		if err = r.valueRepositories.Float.BatchDeleteTx(ctx, &selection_condition.SelectionCondition{
			Where: &float_value.FloatValue{
				EntityID: id,
			},
		}, tx); err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		if err = r.valueRepositories.Date.BatchDeleteTx(ctx, &selection_condition.SelectionCondition{
			Where: &date_value.DateValue{
				EntityID: id,
			},
		}, tx); err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		if err = r.valueRepositories.Time.BatchDeleteTx(ctx, &selection_condition.SelectionCondition{
			Where: &time_value.TimeValue{
				EntityID: id,
			},
		}, tx); err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		if err = r.valueRepositories.Text.BatchDeleteTx(ctx, &selection_condition.SelectionCondition{
			Where: &text_value.TextValue{
				EntityID: id,
			},
		}, tx); err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		if err = tx.Delete(&domain_entity.Entity{}, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.ErrNotFound
			}
		}
		return err
	})
}

//func (r *EntityRepository) textValues2PropertiesValuesMap(entity *domain_entity.Entity) {
//	for _, textValue := range entity.TextValues {
//		entity.PropertiesValuesMap[textValue.PropertyID] = textValue.Value
//	}
//}
//
//// (!) Только если свой-ства устанавливались методом Service.SetValueForProperty() или после изменения PropertiesValuesMap был запущен метод Service.EntityInit(),
//// т.е. состав свойств в PropertiesValues - актуальный
//func (r *EntityRepository) getTextValuesFromPropertiesValuesMap(entity *domain_entity.Entity) (textPropertiesIDs []uint, textValuesMap map[uint]string, err error) {
//	textValuesMap = make(map[uint]string)
//	textPropertiesIDs = make([]uint, 0)
//
//	for propertyID, val := range entity.PropertiesValuesMap {
//		if _, ok := entity.RelationsValues[propertyID]; ok || entity.PropertiesValues[propertyID].Property.PropertyTypeID != property_type.IDText {
//			continue
//		}
//
//		value, ok := val.(string)
//		if !ok {
//			return nil, nil, errors.Errorf("Can not cast string value from PropertiesValuesMap to string. Value = %v.", val)
//		}
//		textValuesMap[propertyID] = value
//		textPropertiesIDs = append(textPropertiesIDs, propertyID)
//	}
//	return textPropertiesIDs, textValuesMap, nil
//}

func (r *EntityRepository) preload(db *gorm.DB, langID uint) *gorm.DB {
	return db.Preload("BoolValues").Preload("IntValues").Preload("FloatValues").Preload("DateValues").Preload("TimeValues").Preload("TextValues", r.db.DB().Model(&text_value.TextValue{}).Where(&text_value.TextValue{LangID: langID}))
}
