package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/property_group"

	"github.com/jinzhu/gorm"
	"github.com/yaruz/app/internal/pkg/apperror"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// PropertyGroupRepository is a repository for the model entity
type PropertyGroupRepository struct {
	repository
}

var _ property_group.Repository = (*PropertyGroupRepository)(nil)

// New creates a new PropertyGroupRepository
func NewPropertyGroupRepository(repository *repository) (*PropertyGroupRepository, error) {
	r := &PropertyGroupRepository{repository: *repository}
	r.autoMigrate()
	return r, nil
}

func (r *PropertyGroupRepository) autoMigrate() {
	if r.db.IsAutoMigrate() {
		r.db.DB().AutoMigrate(&property_group.PropertyGroup{})
	}
}

// Get reads the album with the specified ID from the database.
func (r *PropertyGroupRepository) Get(ctx context.Context, id uint) (*property_group.PropertyGroup, error) {
	entity := &property_group.PropertyGroup{}

	err := r.DB().First(entity, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, yaruzerror.ErrNotFound
		}
	}
	return entity, err
}

func (r *PropertyGroupRepository) First(ctx context.Context, entity *property_group.PropertyGroup) (*property_group.PropertyGroup, error) {
	err := r.DB().Where(entity).First(entity).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, yaruzerror.ErrNotFound
		}
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *PropertyGroupRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property_group.PropertyGroup, error) {
	items := []property_group.PropertyGroup{}
	db := minipkg_gorm.Conditions(r.DB().Model(&property_group.PropertyGroup{}), cond)
	if db.Error != nil {
		return nil, db.Error
	}

	err := db.Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return items, yaruzerror.ErrNotFound
		}
	}
	return items, err
}

func (r *PropertyGroupRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (uint, error) {
	var count uint
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.DB().Model(&property_group.PropertyGroup{}), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	err := db.Count(&count).Error
	return count, err
}

// Create saves a new record in the database.
func (r *PropertyGroupRepository) Create(ctx context.Context, entity *property_group.PropertyGroup) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *PropertyGroupRepository) Update(ctx context.Context, entity *property_group.PropertyGroup) error {

	if r.db.DB().NewRecord(entity) {
		return errors.New("entity is new")
	}
	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *PropertyGroupRepository) Save(ctx context.Context, entity *property_group.PropertyGroup) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *PropertyGroupRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(&property_group.PropertyGroup{}, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return apperror.ErrNotFound
		}
	}
	return err
}
