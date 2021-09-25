package gorm

import (
	"context"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"
	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// EntityRepository is a repository for the model entity
type EntityRepository struct {
	repository
	textValueRepository text_value.Repository
}

var _ entity.Repository = (*EntityRepository)(nil)

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
		r.db.DB().AutoMigrate(&entity.Entity{})
	}
}

// Get reads the album with the specified ID from the database.
func (r *EntityRepository) Get(ctx context.Context, id uint, langID uint) (*entity.Entity, error) {
	entity := &entity.Entity{}

	err := r.textValueJoin(r.db.DB(), langID).First(entity, id).Error
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

func (r *EntityRepository) First(ctx context.Context, entity *entity.Entity, langID uint) (*entity.Entity, error) {
	err := r.textValueJoin(r.db.DB(), langID).Where(entity).First(entity).Error
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
func (r *EntityRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]entity.Entity, error) {
	items := []entity.Entity{}
	db := minipkg_gorm.Conditions(r.textValueJoin(r.db.DB(), langID), cond)
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
func (r *EntityRepository) Create(ctx context.Context, entity *entity.Entity, langID uint) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}

	return r.db.DB().Transaction(func(tx *gorm.DB) (err error) {

		if err := r.beforeSaveTx(ctx, entity, langID, tx); err != nil {
			return err
		}

		return tx.Create(entity).Error
	})
}

// Update saves a changed Maintenance record in the database.
func (r *EntityRepository) Update(ctx context.Context, entity *entity.Entity, langID uint) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}

	return r.Save(ctx, entity, langID)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *EntityRepository) Save(ctx context.Context, entity *entity.Entity, langID uint) error {

	return r.db.DB().Transaction(func(tx *gorm.DB) (err error) {

		if err := r.beforeSaveTx(ctx, entity, langID, tx); err != nil {
			return err
		}

		return tx.Save(entity).Error
	})
}

func (r *EntityRepository) afterFind(ctx context.Context, entity *entity.Entity) error {

	if err := entity.AfterFind(); err != nil {
		return err
	}
	r.textValues2PropertiesValuesMap(entity)
	return nil
}

func (r *EntityRepository) beforeSaveTx(ctx context.Context, entity *entity.Entity, langID uint, tx *gorm.DB) error {
	propertiesIDs, err := r.saveTextValuesFromPropertiesValuesMapTx(ctx, entity, langID, tx)
	if err != nil {
		return err
	}
	r.resetTextValuesInPropertiesValuesMap(ctx, entity, propertiesIDs)

	return entity.BeforeSave()
}

func (r *EntityRepository) resetTextValuesInPropertiesValuesMap(ctx context.Context, entity *entity.Entity, propertiesIDs []uint) {
	for _, propertyID := range propertiesIDs {
		entity.PropertiesValuesMap[propertyID] = 0
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

		if err = tx.Delete(&entity.Entity{}, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.ErrNotFound
			}
		}
		return err
	})
}

func (r *EntityRepository) textValues2PropertiesValuesMap(entity *entity.Entity) {
	for _, textValue := range entity.TextValues {
		entity.PropertiesValuesMap[textValue.PropertyID] = textValue.Value
	}
}

// (!) Только если свой-ства устанавливались методом Service.EntitySetPropertyValue() или после изменения PropertiesValuesMap был запущен метод Service.EntityInit(),
// т.е. состав свойств в PropertiesValues - актуальный
func (r *EntityRepository) saveTextValuesFromPropertiesValuesMapTx(ctx context.Context, entity *entity.Entity, langID uint, tx *gorm.DB) ([]uint, error) {
	newTextValues := make(map[uint]string)
	textPropertiesIDs := make([]uint, 0)

	for propertyID, val := range entity.PropertiesValuesMap {
		if entity.PropertiesValues[propertyID].Property.PropertyTypeID != property_type.IDText {
			continue
		}

		value, ok := val.(string)
		if !ok {
			return nil, errors.Errorf("Can not cast string value from PropertiesValuesMap to string. Value = %v.", val)
		}
		newTextValues[propertyID] = value
		textPropertiesIDs = append(textPropertiesIDs, propertyID)
	}
	return textPropertiesIDs, r.textValueRepository.BatchSaveChangesTx(ctx, entity.ID, langID, newTextValues, tx)
}

func (r *EntityRepository) textValueJoin(db *gorm.DB, langID uint) *gorm.DB {
	//return db.Joins("left join text_value on entity.id = text_value.entity_id and lang_id = ?", langID)
	return db.Joins("TextValues", r.db.DB().Where(&text_value.TextValue{LangID: langID}))
}
