package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yaruz_data/domain/property_type"

	"github.com/jinzhu/gorm"
	"github.com/yaruz/app/internal/pkg/apperror"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yaruz_data/pkgerror"
)

// PropertyTypeRepository is a repository for the model entity
type PropertyTypeRepository struct {
	repository
}

var _ property_type.Repository = (*PropertyTypeRepository)(nil)

// New creates a new PropertyTypeRepository
func NewPropertyTypeRepository(repository *repository) (*PropertyTypeRepository, error) {
	return &PropertyTypeRepository{repository: *repository}, nil
}

func (r *PropertyTypeRepository) autoMigrate() {
	if r.db.IsAutoMigrate() {
		r.db.DB().AutoMigrate(&property_type.PropertyType{})
	}
}

// Get reads the album with the specified ID from the database.
func (r *PropertyTypeRepository) Get(ctx context.Context, id uint) (*property_type.PropertyType, error) {
	entity := &property_type.PropertyType{}

	err := r.DB().First(entity, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, pkgerror.ErrNotFound
		}
	}
	return entity, err
}

func (r *PropertyTypeRepository) First(ctx context.Context, entity *property_type.PropertyType) (*property_type.PropertyType, error) {
	err := r.DB().Where(entity).First(entity).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, pkgerror.ErrNotFound
		}
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *PropertyTypeRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property_type.PropertyType, error) {
	items := []property_type.PropertyType{}
	db := minipkg_gorm.Conditions(r.DB().Model(&property_type.PropertyType{}), cond)
	if db.Error != nil {
		return nil, db.Error
	}

	err := db.Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return items, pkgerror.ErrNotFound
		}
	}
	return items, err
}

func (r *PropertyTypeRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (uint, error) {
	var count uint
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.DB().Model(&property_type.PropertyType{}), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	err := db.Count(&count).Error
	return count, err
}

// Create saves a new record in the database.
func (r *PropertyTypeRepository) Create(ctx context.Context, entity *property_type.PropertyType) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *PropertyTypeRepository) Update(ctx context.Context, entity *property_type.PropertyType) error {

	if r.db.DB().NewRecord(entity) {
		return errors.New("entity is new")
	}
	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *PropertyTypeRepository) Save(ctx context.Context, entity *property_type.PropertyType) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *PropertyTypeRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(&property_type.PropertyType{}, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return apperror.ErrNotFound
		}
	}
	return err
}
