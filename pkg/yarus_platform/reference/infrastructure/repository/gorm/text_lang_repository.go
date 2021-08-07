package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_lang"

	"github.com/yaruz/app/internal/pkg/apperror"
	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// TextLangRepository is a repository for the model entity
type TextLangRepository struct {
	repository
}

var _ text_lang.Repository = (*TextLangRepository)(nil)

// New creates a new TextLangRepository
func NewTextLangRepository(repository *repository) (*TextLangRepository, error) {
	return &TextLangRepository{repository: *repository}, nil
}

// Get reads the album with the specified ID from the database.
func (r *TextLangRepository) Get(ctx context.Context, id uint) (*text_lang.TextLang, error) {
	entity := &text_lang.TextLang{}

	err := r.db.DB().First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

func (r *TextLangRepository) First(ctx context.Context, entity *text_lang.TextLang) (*text_lang.TextLang, error) {
	err := r.db.DB().Where(entity).First(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *TextLangRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]text_lang.TextLang, error) {
	items := []text_lang.TextLang{}
	db := minipkg_gorm.Conditions(r.db.DB(), cond)
	if db.Error != nil {
		return nil, db.Error
	}

	err := db.Model(r.model).Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return items, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return items, err
}

func (r *TextLangRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
	var count int64
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.db.DB(), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	return count, db.Model(r.model).Count(&count).Error
}

// Create saves a new record in the database.
func (r *TextLangRepository) Create(ctx context.Context, entity *text_lang.TextLang) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}

// Update saves a changed record in the database.
func (r *TextLangRepository) Update(ctx context.Context, entity *text_lang.TextLang) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}
	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *TextLangRepository) Save(ctx context.Context, entity *text_lang.TextLang) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a record in the database.
func (r *TextLangRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(r.model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}
