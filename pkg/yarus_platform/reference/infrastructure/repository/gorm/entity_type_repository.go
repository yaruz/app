package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type2property"

	"github.com/yaruz/app/internal/pkg/apperror"
	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// EntityTypeRepository is a repository for the model entity
type EntityTypeRepository struct {
	repository
	entityType2PropertyRepository entity_type2property.Repository
}

var _ entity_type.Repository = (*EntityTypeRepository)(nil)

// New creates a new EntityTypeRepository
func NewEntityTypeRepository(repository *repository, entityType2PropertyRepository *entity_type2property.Repository) (*EntityTypeRepository, error) {
	return &EntityTypeRepository{repository: *repository, entityType2PropertyRepository: *entityType2PropertyRepository}, nil
}

// Get reads the album with the specified ID from the database.
func (r *EntityTypeRepository) Get(ctx context.Context, id uint) (*entity_type.EntityType, error) {
	entity := &entity_type.EntityType{}

	err := r.db.DB().First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

func (r *EntityTypeRepository) First(ctx context.Context, entity *entity_type.EntityType) (*entity_type.EntityType, error) {
	err := r.db.DB().Where(entity).First(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *EntityTypeRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]entity_type.EntityType, error) {
	items := []entity_type.EntityType{}
	db := minipkg_gorm.Conditions(r.db.DB(), cond)
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

func (r *EntityTypeRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
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
func (r *EntityTypeRepository) Create(ctx context.Context, entity *entity_type.EntityType) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *EntityTypeRepository) Update(ctx context.Context, entity *entity_type.EntityType) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}
	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *EntityTypeRepository) Save(ctx context.Context, entity *entity_type.EntityType) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *EntityTypeRepository) Delete(ctx context.Context, id uint) error {
	return r.db.DB().Transaction(func(tx *gorm.DB) error {
		if err := r.unbindAllPropertyTx(ctx, tx, id); err != nil {
			return err
		}

		err := tx.Delete(r.model, id).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.ErrNotFound
			}
		}

		return err
	})
}

func (r *EntityTypeRepository) BindProperty(ctx context.Context, id uint, propertyID uint) error {
	// todo: хорошо бы сделать валидацию на тип свойства - все, кроме связей
	return r.entityType2PropertyRepository.Create(ctx, &entity_type2property.EntityType2Property{
		EntityTypeID: id,
		PropertyID:   propertyID,
	})
}

func (r *EntityTypeRepository) UnbindProperty(ctx context.Context, id uint, propertyID uint) error {
	return r.entityType2PropertyRepository.Delete(ctx, &entity_type2property.EntityType2Property{
		EntityTypeID: id,
		PropertyID:   propertyID,
	})
}

func (r *EntityTypeRepository) unbindAllPropertyTx(ctx context.Context, tx *gorm.DB, id uint) error {
	return r.entityType2PropertyRepository.DeleteTx(ctx, tx, &entity_type2property.EntityType2Property{
		EntityTypeID: id,
	})
}

func (r *EntityTypeRepository) UnbindAllProperty(ctx context.Context, id uint) error {
	return r.entityType2PropertyRepository.Delete(ctx, &entity_type2property.EntityType2Property{
		EntityTypeID: id,
	})
}
