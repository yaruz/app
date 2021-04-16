package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_value"

	"github.com/yaruz/app/internal/pkg/apperror"

	"github.com/jinzhu/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// TextValueRepository is a repository for the model entity
type TextValueRepository struct {
	repository
}

var _ text_value.Repository = (*TextValueRepository)(nil)

// New creates a new TextValueRepository
func NewTextValueRepository(repository *repository) (*TextValueRepository, error) {
	r := &TextValueRepository{repository: *repository}
	r.autoMigrate()
	return r, nil
}

func (r *TextValueRepository) autoMigrate() {
	if r.db.IsAutoMigrate() {
		r.db.DB().AutoMigrate(&text_value.TextValue{})
	}
}

// Get reads the album with the specified ID from the database.
func (r *TextValueRepository) Get(ctx context.Context, id uint) (*text_value.TextValue, error) {
	entity := &text_value.TextValue{}

	err := r.DB().First(entity, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, yaruzerror.ErrNotFound
		}
	}

	return entity, err
}

func (r *TextValueRepository) First(ctx context.Context, entity *text_value.TextValue) (*text_value.TextValue, error) {
	err := r.DB().Where(entity).First(entity).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, yaruzerror.ErrNotFound
		}
	}

	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *TextValueRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]text_value.TextValue, error) {
	items := []text_value.TextValue{}
	db := minipkg_gorm.Conditions(r.DB().Model(&text_value.TextValue{}), cond)
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

func (r *TextValueRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (uint, error) {
	var count uint
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.DB().Model(&text_value.TextValue{}), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	err := db.Count(&count).Error
	return count, err
}

// Create saves a new record in the database.
func (r *TextValueRepository) Create(ctx context.Context, entity *text_value.TextValue) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}

	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *TextValueRepository) Update(ctx context.Context, entity *text_value.TextValue) error {

	if r.db.DB().NewRecord(entity) {
		return errors.New("entity is new")
	}

	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *TextValueRepository) Save(ctx context.Context, entity *text_value.TextValue) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *TextValueRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(&text_value.TextValue{}, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return apperror.ErrNotFound
		}
	}
	return err
}
