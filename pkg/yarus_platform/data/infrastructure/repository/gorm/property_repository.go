package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/property"

	"github.com/jinzhu/gorm"
	"github.com/yaruz/app/internal/pkg/apperror"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// PropertyRepository is a repository for the model entity
type PropertyRepository struct {
	repository
}

var _ property.Repository = (*PropertyRepository)(nil)

// New creates a new PropertyRepository
func NewPropertyRepository(repository *repository) (*PropertyRepository, error) {
	r := &PropertyRepository{repository: *repository}
	r.autoMigrate()
	return r, nil
}

func (r *PropertyRepository) autoMigrate() {
	if r.db.IsAutoMigrate() {
		r.db.DB().AutoMigrate(&property.Property{})
	}
}

// Get reads the album with the specified ID from the database.
func (r *PropertyRepository) Get(ctx context.Context, id uint) (*property.Property, error) {
	entity := &property.Property{}

	err := r.DB().First(entity, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, yaruzerror.ErrNotFound
		}
	}

	if err = entity.AfterFind(); err != nil {
		return nil, err
	}

	return entity, err
}

func (r *PropertyRepository) First(ctx context.Context, entity *property.Property) (*property.Property, error) {
	err := r.DB().Where(entity).First(entity).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, yaruzerror.ErrNotFound
		}
	}

	if err = entity.AfterFind(); err != nil {
		return nil, err
	}

	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *PropertyRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property.Property, error) {
	items := []property.Property{}
	db := minipkg_gorm.Conditions(r.DB().Model(&property.Property{}), cond)
	if db.Error != nil {
		return nil, db.Error
	}

	err := db.Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return items, yaruzerror.ErrNotFound
		}
	}

	for _, entity := range items {
		if err = entity.AfterFind(); err != nil {
			return nil, err
		}
	}

	return items, err
}

func (r *PropertyRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (uint, error) {
	var count uint
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.DB().Model(&property.Property{}), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	err := db.Count(&count).Error
	return count, err
}

// Create saves a new record in the database.
func (r *PropertyRepository) Create(ctx context.Context, entity *property.Property) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}

	if err := entity.BeforeSave(); err != nil {
		return err
	}

	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *PropertyRepository) Update(ctx context.Context, entity *property.Property) error {

	if r.db.DB().NewRecord(entity) {
		return errors.New("entity is new")
	}

	if err := entity.BeforeSave(); err != nil {
		return err
	}

	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *PropertyRepository) Save(ctx context.Context, entity *property.Property) error {

	if err := entity.BeforeSave(); err != nil {
		return err
	}

	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *PropertyRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(&property.Property{}, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return apperror.ErrNotFound
		}
	}
	return err
}
