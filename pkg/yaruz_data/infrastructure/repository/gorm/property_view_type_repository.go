package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yaruz_data/domain/property_view_type"

	"github.com/jinzhu/gorm"
	"github.com/yaruz/app/internal/pkg/apperror"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yaruz_data/pkgerror"
)

// PropertyViewTypeRepository is a repository for the model entity
type PropertyViewTypeRepository struct {
	repository
}

var _ property_view_type.Repository = (*PropertyViewTypeRepository)(nil)

// New creates a new PropertyViewTypeRepository
func NewPropertyViewTypeRepository(repository *repository) (*PropertyViewTypeRepository, error) {
	return &PropertyViewTypeRepository{repository: *repository}, nil
}

func (r *PropertyViewTypeRepository) autoMigrate() {
	if r.db.IsAutoMigrate() {
		r.db.DB().AutoMigrate(&property_view_type.PropertyViewType{})
	}
}

// Get reads the album with the specified ID from the database.
func (r *PropertyViewTypeRepository) Get(ctx context.Context, id uint) (*property_view_type.PropertyViewType, error) {
	entity := &property_view_type.PropertyViewType{}

	err := r.DB().First(entity, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, pkgerror.ErrNotFound
		}
	}
	return entity, err
}

func (r *PropertyViewTypeRepository) First(ctx context.Context, entity *property_view_type.PropertyViewType) (*property_view_type.PropertyViewType, error) {
	err := r.DB().Where(entity).First(entity).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, pkgerror.ErrNotFound
		}
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *PropertyViewTypeRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property_view_type.PropertyViewType, error) {
	items := []property_view_type.PropertyViewType{}
	db := minipkg_gorm.Conditions(r.DB().Model(&property_view_type.PropertyViewType{}), cond)
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

func (r *PropertyViewTypeRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (uint, error) {
	var count uint
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.DB().Model(&property_view_type.PropertyViewType{}), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	err := db.Count(&count).Error
	return count, err
}

// Create saves a new record in the database.
func (r *PropertyViewTypeRepository) Create(ctx context.Context, entity *property_view_type.PropertyViewType) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *PropertyViewTypeRepository) Update(ctx context.Context, entity *property_view_type.PropertyViewType) error {

	if r.db.DB().NewRecord(entity) {
		return errors.New("entity is new")
	}
	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *PropertyViewTypeRepository) Save(ctx context.Context, entity *property_view_type.PropertyViewType) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *PropertyViewTypeRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(&property_view_type.PropertyViewType{}, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return apperror.ErrNotFound
		}
	}
	return err
}
