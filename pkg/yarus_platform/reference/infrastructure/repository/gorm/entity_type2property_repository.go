package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type2property"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"

	"github.com/yaruz/app/internal/pkg/apperror"
	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"
)

// EntityType2PropertyRepository is a repository for the EntityType2Property entity
type EntityType2PropertyRepository struct {
	repository
}

var _ entity_type2property.Repository = (*EntityType2PropertyRepository)(nil)

// New creates a new EntityType2PropertyRepository
func NewEntityType2PropertyRepository(repository *repository) (*EntityType2PropertyRepository, error) {
	return &EntityType2PropertyRepository{repository: *repository}, nil
}

// Get reads the album with the specified ID from the database.
func (r *EntityType2PropertyRepository) Get(ctx context.Context, id uint) (*entity_type2property.EntityType2Property, error) {
	entity := &entity_type2property.EntityType2Property{}

	err := r.db.DB().First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruserror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

func (r *EntityType2PropertyRepository) First(ctx context.Context, entity *entity_type2property.EntityType2Property) (*entity_type2property.EntityType2Property, error) {
	err := r.db.DB().Where(entity).First(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruserror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *EntityType2PropertyRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]entity_type2property.EntityType2Property, error) {
	return r.QueryTx(ctx, r.db.DB(), cond)
}

func (r *EntityType2PropertyRepository) QueryTx(ctx context.Context, tx *gorm.DB, cond *selection_condition.SelectionCondition) ([]entity_type2property.EntityType2Property, error) {
	items := []entity_type2property.EntityType2Property{}
	db := minipkg_gorm.Conditions(r.db.GormTx(tx), cond)
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

func (r *EntityType2PropertyRepository) QueryWithEntityType(ctx context.Context, cond *selection_condition.SelectionCondition) ([]entity_type2property.EntityType2Property, error) {
	items := []entity_type2property.EntityType2Property{}
	db := minipkg_gorm.Conditions(r.db.DB(), cond)
	if db.Error != nil {
		return nil, db.Error
	}

	err := r.joins(db).Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return items, yaruserror.ErrNotFound
		}
		return nil, err
	}
	return items, err
}

func (r *EntityType2PropertyRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
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
func (r *EntityType2PropertyRepository) Create(ctx context.Context, entity *entity_type2property.EntityType2Property) error {
	return r.db.DB().Create(entity).Error
}
func (r *EntityType2PropertyRepository) CreateTx(ctx context.Context, tx *gorm.DB, entity *entity_type2property.EntityType2Property) error {
	return r.db.GormTx(tx).Create(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *EntityType2PropertyRepository) Delete(ctx context.Context, entity *entity_type2property.EntityType2Property) error {

	err := r.db.DB().Delete(r.model, entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}

func (r *EntityType2PropertyRepository) DeleteTx(ctx context.Context, tx *gorm.DB, entity *entity_type2property.EntityType2Property) error {

	err := r.db.GormTx(tx).Delete(r.model, entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}

func (r *EntityType2PropertyRepository) joins(db *gorm.DB) *gorm.DB {
	return db.Joins("EntityType")
}
