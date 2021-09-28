package gorm

import (
	"context"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/yaruz/app/internal/pkg/apperror"
	domain_entity "github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"
	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// EntityRepository is a repository for the model entity
type EntityRepository struct {
	repository
	textValueRepository text_value.Repository
}

var _ domain_entity.Repository = (*EntityRepository)(nil)

// New creates a new EntityRepository
func NewEntityRepository(repository *repository, textValueRepository text_value.Repository) (*EntityRepository, error) {
	r := &EntityRepository{
		repository:          *repository,
		textValueRepository: textValueRepository,
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

	err := r.textValuePreload(r.db.DB(), langID).First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}

	if err = r.afterFind(ctx, entity); err != nil {
		return nil, err
	}

	return entity, err
}

func (r *EntityRepository) First(ctx context.Context, entity *domain_entity.Entity, langID uint) (*domain_entity.Entity, error) {
	err := r.textValuePreload(r.db.DB(), langID).Where(entity).First(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}

	if err = r.afterFind(ctx, entity); err != nil {
		return nil, err
	}

	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *EntityRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]domain_entity.Entity, error) {
	items := []domain_entity.Entity{}
	db := minipkg_gorm.Conditions(r.textValuePreload(r.db.DB(), langID), cond)
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

	for i := range items {
		if err = r.afterFind(ctx, &items[i]); err != nil {
			return nil, err
		}
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

	textValuesMap, err := r.beforeSave(ctx, entity)
	if err != nil {
		return err
	}

	return r.db.DB().Transaction(func(tx *gorm.DB) (err error) {

		if err := tx.Create(entity).Error; err != nil {
			return err
		}

		return r.afterSaveTx(ctx, entity, textValuesMap, langID, tx)
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

	textValuesMap, err := r.beforeSave(ctx, entity)
	if err != nil {
		return err
	}

	return tx.Transaction(func(tx *gorm.DB) (err error) {

		if err := tx.Save(entity).Error; err != nil {
			return err
		}

		return r.afterSaveTx(ctx, entity, textValuesMap, langID, tx)
	})
}

func (r *EntityRepository) afterFind(ctx context.Context, entity *domain_entity.Entity) error {

	if err := entity.AfterFind(); err != nil {
		return err
	}
	r.textValues2PropertiesValuesMap(entity)
	return nil
}

func (r *EntityRepository) beforeSave(ctx context.Context, entity *domain_entity.Entity) (textValuesMap map[uint]string, err error) {
	propertiesIDs, textValuesMap, err := r.getTextValuesFromPropertiesValuesMap(entity)
	if err != nil {
		return nil, err
	}
	r.resetTextValuesInPropertiesValuesMap(ctx, entity, propertiesIDs)

	return textValuesMap, entity.BeforeSave()
}

func (r *EntityRepository) afterSaveTx(ctx context.Context, entity *domain_entity.Entity, textValuesMap map[uint]string, langID uint, tx *gorm.DB) error {
	for propertyID, value := range textValuesMap {
		entity.PropertiesValuesMap[propertyID] = value
	}
	return r.textValueRepository.BatchSaveChangesTx(ctx, entity.ID, langID, textValuesMap, tx)
}

func (r *EntityRepository) resetTextValuesInPropertiesValuesMap(ctx context.Context, entity *domain_entity.Entity, propertiesIDs []uint) {
	for _, propertyID := range propertiesIDs {
		entity.PropertiesValuesMap[propertyID] = ""
	}
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *EntityRepository) Delete(ctx context.Context, id uint) error {

	return r.db.DB().Transaction(func(tx *gorm.DB) (err error) {
		if err = r.textValueRepository.BatchDeleteTx(ctx, &selection_condition.SelectionCondition{
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

func (r *EntityRepository) textValues2PropertiesValuesMap(entity *domain_entity.Entity) {
	for _, textValue := range entity.TextValues {
		entity.PropertiesValuesMap[textValue.PropertyID] = textValue.Value
	}
}

// (!) Только если свой-ства устанавливались методом Service.EntitySetPropertyValue() или после изменения PropertiesValuesMap был запущен метод Service.EntityInit(),
// т.е. состав свойств в PropertiesValues - актуальный
func (r *EntityRepository) getTextValuesFromPropertiesValuesMap(entity *domain_entity.Entity) (textPropertiesIDs []uint, textValuesMap map[uint]string, err error) {
	textValuesMap = make(map[uint]string)
	textPropertiesIDs = make([]uint, 0)

	for propertyID, val := range entity.PropertiesValuesMap {
		if entity.PropertiesValues[propertyID].Property.PropertyTypeID != property_type.IDText {
			continue
		}

		value, ok := val.(string)
		if !ok {
			return nil, nil, errors.Errorf("Can not cast string value from PropertiesValuesMap to string. Value = %v.", val)
		}
		textValuesMap[propertyID] = value
		textPropertiesIDs = append(textPropertiesIDs, propertyID)
	}
	return textPropertiesIDs, textValuesMap, nil
}

func (r *EntityRepository) textValuePreload(db *gorm.DB, langID uint) *gorm.DB {
	return db.Preload("TextValues", r.db.DB().Model(&text_value.TextValue{}).Where(&text_value.TextValue{LangID: langID}))
}
