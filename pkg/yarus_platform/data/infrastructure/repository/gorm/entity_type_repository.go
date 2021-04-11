package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/pkg/yarus_platform/data/domain/entity_type"

	"github.com/jinzhu/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// EntityTypeRepository is a repository for the model entity
type EntityTypeRepository struct {
	repository
}

var _ entity_type.Repository = (*EntityTypeRepository)(nil)

// New creates a new EntityTypeRepository
func NewEntityTypeRepository(repository *repository) (*EntityTypeRepository, error) {
	return &EntityTypeRepository{repository: *repository}, nil
}

func (r *EntityTypeRepository) autoMigrate() {
	if r.db.IsAutoMigrate() {
		r.db.DB().AutoMigrate(&entity_type.EntityType{})
	}
}

// Get reads the album with the specified ID from the database.
func (r *EntityTypeRepository) Get(ctx context.Context, id uint) (*entity_type.EntityType, error) {
	entity := &entity_type.EntityType{}

	err := r.DB().First(entity, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, yaruzerror.ErrNotFound
		}
	}
	return entity, err
}

func (r *EntityTypeRepository) First(ctx context.Context, entity *entity_type.EntityType) (*entity_type.EntityType, error) {
	err := r.DB().Where(entity).First(entity).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, yaruzerror.ErrNotFound
		}
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *EntityTypeRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]entity_type.EntityType, error) {
	items := []entity_type.EntityType{}
	db := minipkg_gorm.Conditions(r.DB().Model(&entity_type.EntityType{}), cond)
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

func (r *EntityTypeRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (uint, error) {
	var count uint
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.DB().Model(&entity_type.EntityType{}), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	err := db.Count(&count).Error
	return count, err
}

// Create saves a new record in the database.
func (r *EntityTypeRepository) Create(ctx context.Context, entity *entity_type.EntityType) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *EntityTypeRepository) Update(ctx context.Context, entity *entity_type.EntityType) error {

	if r.db.DB().NewRecord(entity) {
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

	err := r.db.DB().Delete(&entity_type.EntityType{}, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return apperror.ErrNotFound
		}
	}
	return err
}
