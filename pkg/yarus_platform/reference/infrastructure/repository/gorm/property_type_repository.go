package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"

	"github.com/yaruz/app/internal/pkg/apperror"
	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
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

// Get reads the album with the specified ID from the database.
func (r *PropertyTypeRepository) Get(ctx context.Context, id uint) (*property_type.PropertyType, error) {
	entity := &property_type.PropertyType{}

	err := r.DB().First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

func (r *PropertyTypeRepository) First(ctx context.Context, entity *property_type.PropertyType) (*property_type.PropertyType, error) {
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
func (r *PropertyTypeRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property_type.PropertyType, error) {
	items := []property_type.PropertyType{}
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

func (r *PropertyTypeRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
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
func (r *PropertyTypeRepository) Create(ctx context.Context, entity *property_type.PropertyType) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}

// Update saves a changed record in the database.
func (r *PropertyTypeRepository) Update(ctx context.Context, entity *property_type.PropertyType) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}
	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *PropertyTypeRepository) Save(ctx context.Context, entity *property_type.PropertyType) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a record in the database.
func (r *PropertyTypeRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(&property_type.PropertyType{}, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}
