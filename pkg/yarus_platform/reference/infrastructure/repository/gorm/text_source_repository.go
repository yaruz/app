package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_source"

	"github.com/yaruz/app/internal/pkg/apperror"

	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// TextSourceRepository is a repository for the model entity
type TextSourceRepository struct {
	repository
}

var _ text_source.Repository = (*TextSourceRepository)(nil)

// New creates a new TextSourceRepository
func NewTextSourceRepository(repository *repository) (*TextSourceRepository, error) {
	return &TextSourceRepository{repository: *repository}, nil
}

// Get reads the album with the specified ID from the database.
func (r *TextSourceRepository) Get(ctx context.Context, id uint) (*text_source.TextSource, error) {
	entity := &text_source.TextSource{}

	err := r.DB().First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

func (r *TextSourceRepository) First(ctx context.Context, entity *text_source.TextSource) (*text_source.TextSource, error) {
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
func (r *TextSourceRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]text_source.TextSource, error) {
	items := []text_source.TextSource{}
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

func (r *TextSourceRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
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
func (r *TextSourceRepository) Create(ctx context.Context, entity *text_source.TextSource) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}

	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *TextSourceRepository) Update(ctx context.Context, entity *text_source.TextSource) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}

	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *TextSourceRepository) Save(ctx context.Context, entity *text_source.TextSource) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *TextSourceRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(r.model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}
