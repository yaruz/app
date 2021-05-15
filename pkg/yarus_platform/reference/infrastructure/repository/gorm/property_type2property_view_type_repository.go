package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type2property_view_type"

	"github.com/yaruz/app/internal/pkg/apperror"
	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// PropertyType2PropertyViewTypeRepository is a repository for the model entity
type PropertyType2PropertyViewTypeRepository struct {
	repository
}

var _ property_type2property_view_type.Repository = (*PropertyType2PropertyViewTypeRepository)(nil)

// New creates a new PropertyType2PropertyViewTypeRepository
func NewPropertyType2PropertyViewTypeRepository(repository *repository) (*PropertyType2PropertyViewTypeRepository, error) {
	return &PropertyType2PropertyViewTypeRepository{repository: *repository}, nil
}

// Get reads the album with the specified ID from the database.
func (r *PropertyType2PropertyViewTypeRepository) Get(ctx context.Context, id uint) (*property_type2property_view_type.PropertyType2PropertyViewType, error) {
	entity := &property_type2property_view_type.PropertyType2PropertyViewType{}

	err := r.DB().First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

func (r *PropertyType2PropertyViewTypeRepository) First(ctx context.Context, entity *property_type2property_view_type.PropertyType2PropertyViewType) (*property_type2property_view_type.PropertyType2PropertyViewType, error) {
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
func (r *PropertyType2PropertyViewTypeRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property_type2property_view_type.PropertyType2PropertyViewType, error) {
	items := []property_type2property_view_type.PropertyType2PropertyViewType{}
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

func (r *PropertyType2PropertyViewTypeRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
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
func (r *PropertyType2PropertyViewTypeRepository) Create(ctx context.Context, entity *property_type2property_view_type.PropertyType2PropertyViewType) error {
	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *PropertyType2PropertyViewTypeRepository) Update(ctx context.Context, entity *property_type2property_view_type.PropertyType2PropertyViewType) error {
	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *PropertyType2PropertyViewTypeRepository) Save(ctx context.Context, entity *property_type2property_view_type.PropertyType2PropertyViewType) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *PropertyType2PropertyViewTypeRepository) Delete(ctx context.Context, propertyTypeID uint, propertyViewTypeID uint) error {

	err := r.db.DB().Delete(r.model, &property_type2property_view_type.PropertyType2PropertyViewType{
		PropertyTypeID:     propertyTypeID,
		PropertyViewTypeID: propertyViewTypeID,
	}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}
