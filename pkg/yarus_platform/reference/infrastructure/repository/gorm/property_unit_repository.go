package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_unit"

	"github.com/jinzhu/gorm"
	"github.com/yaruz/app/internal/pkg/apperror"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// PropertyUnitRepository is a repository for the model entity
type PropertyUnitRepository struct {
	repository
}

var _ property_unit.Repository = (*PropertyUnitRepository)(nil)

// New creates a new PropertyUnitRepository
func NewPropertyUnitRepository(repository *repository) (*PropertyUnitRepository, error) {
	r := &PropertyUnitRepository{repository: *repository}
	r.autoMigrate()
	return r, nil
}

func (r *PropertyUnitRepository) autoMigrate() {
	if r.db.IsAutoMigrate() {
		r.db.DB().AutoMigrate(&property_unit.PropertyUnit{})
	}
}

// Get reads the album with the specified ID from the database.
func (r *PropertyUnitRepository) Get(ctx context.Context, id uint) (*property_unit.PropertyUnit, error) {
	entity := &property_unit.PropertyUnit{}

	err := r.DB().First(entity, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

func (r *PropertyUnitRepository) First(ctx context.Context, entity *property_unit.PropertyUnit) (*property_unit.PropertyUnit, error) {
	err := r.DB().Where(entity).First(entity).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *PropertyUnitRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property_unit.PropertyUnit, error) {
	items := []property_unit.PropertyUnit{}
	db := minipkg_gorm.Conditions(r.DB().Model(&property_unit.PropertyUnit{}), cond)
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

func (r *PropertyUnitRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (uint, error) {
	var count uint
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.DB().Model(&property_unit.PropertyUnit{}), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	err := db.Count(&count).Error
	return count, err
}

// Create saves a new record in the database.
func (r *PropertyUnitRepository) Create(ctx context.Context, entity *property_unit.PropertyUnit) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *PropertyUnitRepository) Update(ctx context.Context, entity *property_unit.PropertyUnit) error {

	if r.db.DB().NewRecord(entity) {
		return errors.New("entity is new")
	}
	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *PropertyUnitRepository) Save(ctx context.Context, entity *property_unit.PropertyUnit) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *PropertyUnitRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(&property_unit.PropertyUnit{}, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return apperror.ErrNotFound
		}
	}
	return err
}
