package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/data/domain/t_string"

	"github.com/yaruz/app/internal/pkg/apperror"

	"github.com/jinzhu/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// TTextRepository is a repository for the model entity
type TTextRepository struct {
	repository
}

var _ t_string.Repository = (*TTextRepository)(nil)

// New creates a new TTextRepository
func NewTTextRepository(repository *repository) (*TTextRepository, error) {
	return &TTextRepository{repository: *repository}, nil
}

func (r *TTextRepository) autoMigrate() {
	if r.db.IsAutoMigrate() {
		r.db.DB().AutoMigrate(&t_string.TString{})
	}
}

// Get reads the album with the specified ID from the database.
func (r *TTextRepository) Get(ctx context.Context, id uint) (*t_string.TString, error) {
	entity := &t_string.TString{}

	err := r.DB().First(entity, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, yaruzerror.ErrNotFound
		}
	}

	return entity, err
}

func (r *TTextRepository) First(ctx context.Context, entity *t_string.TString) (*t_string.TString, error) {
	err := r.DB().Where(entity).First(entity).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity, yaruzerror.ErrNotFound
		}
	}

	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *TTextRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]t_string.TString, error) {
	items := []t_string.TString{}
	db := minipkg_gorm.Conditions(r.DB().Model(&t_string.TString{}), cond)
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

func (r *TTextRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (uint, error) {
	var count uint
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.DB().Model(&t_string.TString{}), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	err := db.Count(&count).Error
	return count, err
}

// Create saves a new record in the database.
func (r *TTextRepository) Create(ctx context.Context, entity *t_string.TString) error {

	if !r.db.DB().NewRecord(entity) {
		return errors.New("entity is not new")
	}

	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *TTextRepository) Update(ctx context.Context, entity *t_string.TString) error {

	if r.db.DB().NewRecord(entity) {
		return errors.New("entity is new")
	}

	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *TTextRepository) Save(ctx context.Context, entity *t_string.TString) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *TTextRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(&t_string.TString{}, id).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return apperror.ErrNotFound
		}
	}
	return err
}
