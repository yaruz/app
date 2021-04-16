package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/text_source"

	"github.com/yaruz/app/internal/pkg/apperror"

	"github.com/jinzhu/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// TSourceRepository is a repository for the model entity
type TSourceRepository struct {
	repository
}

var _ text_source.Repository = (*TSourceRepository)(nil)

// New creates a new TSourceRepository
func NewTSourceRepository(repository *repository) (*TSourceRepository, error) {
	r := &TSourceRepository{repository: *repository}
	r.autoMigrate()
	return r, nil
}

func (r *TSourceRepository) autoMigrate() {
	if r.db.IsAutoMigrate() {
		r.db.DB().AutoMigrate(&text_source.TextSource{})
	}
}

// Get reads the album with the specified ID from the database.
func (r *TSourceRepository) Get(ctx context.Context, id uint) (*text_source.TextSource, error) {
	entity := &text_source.TextSource{}

	err := r.DB().First(entity, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, yaruzerror.ErrNotFound
		}
	}

	return entity, err
}

func (r *TSourceRepository) First(ctx context.Context, entity *text_source.TextSource) (*text_source.TextSource, error) {
	err := r.DB().Where(entity).First(entity).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, yaruzerror.ErrNotFound
		}
	}

	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *TSourceRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]text_source.TextSource, error) {
	items := []text_source.TextSource{}
	db := minipkg_gorm.Conditions(r.DB().Model(&text_source.TextSource{}), cond)
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

func (r *TSourceRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (uint, error) {
	var count uint
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.DB().Model(&text_source.TextSource{}), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	err := db.Count(&count).Error
	return count, err
}

// Create saves a new record in the database.
func (r *TSourceRepository) Create(ctx context.Context, entity *text_source.TextSource) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}

	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *TSourceRepository) Update(ctx context.Context, entity *text_source.TextSource) error {

	if r.db.DB().NewRecord(entity) {
		return errors.New("entity is new")
	}

	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *TSourceRepository) Save(ctx context.Context, entity *text_source.TextSource) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *TSourceRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(&text_source.TextSource{}, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return apperror.ErrNotFound
		}
	}
	return err
}
