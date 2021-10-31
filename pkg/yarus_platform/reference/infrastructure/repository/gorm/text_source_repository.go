package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_value"

	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_source"

	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
)

// TextSourceRepository is a repository for the model entity
type TextSourceRepository struct {
	repository
	textValueRepository text_value.Repository
}

var _ text_source.Repository = (*TextSourceRepository)(nil)

// New creates a new TextSourceRepository
func NewTextSourceRepository(repository *repository, textValueRepository text_value.Repository) (*TextSourceRepository, error) {
	return &TextSourceRepository{
		repository:          *repository,
		textValueRepository: textValueRepository,
	}, nil
}

// Get reads the album with the specified ID from the database.
func (r *TextSourceRepository) Get(ctx context.Context, id uint) (*text_source.TextSource, error) {
	return r.getTx(ctx, r.db.DB(), id)
}

// TGet reads the album with the specified ID from the database.
func (r *TextSourceRepository) TGet(ctx context.Context, id uint, langID uint) (*text_source.TextSource, error) {
	db := r.joins(r.db.DB(), langID)
	return r.getTx(ctx, db, id)
}

func (r *TextSourceRepository) getTx(ctx context.Context, tx *gorm.DB, id uint) (*text_source.TextSource, error) {
	entity := &text_source.TextSource{}

	err := r.db.GormTx(tx).First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruserror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

func (r *TextSourceRepository) First(ctx context.Context, entity *text_source.TextSource) (*text_source.TextSource, error) {
	err := r.db.DB().Where(entity).First(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruserror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *TextSourceRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]text_source.TextSource, error) {
	items := []text_source.TextSource{}
	db := minipkg_gorm.Conditions(r.db.DB(), cond)
	if db.Error != nil {
		return nil, db.Error
	}

	err := db.Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return items, yaruserror.ErrNotFound
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
	db := minipkg_gorm.Conditions(r.db.DB(), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	return count, db.Model(r.model).Count(&count).Error
}

// Create saves a new record in the database.
func (r *TextSourceRepository) Create(ctx context.Context, entity *text_source.TextSource) error {
	return r.createTx(ctx, r.db.DB(), entity)
}

func (r *TextSourceRepository) createTx(ctx context.Context, tx *gorm.DB, entity *text_source.TextSource) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}

	return r.db.GormTx(tx).Create(entity).Error
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
	return r.DeleteTx(ctx, r.db.DB(), id)
}

func (r *TextSourceRepository) DeleteTx(ctx context.Context, tx *gorm.DB, id uint) (err error) {

	if err = r.textValueRepository.DeleteTx(ctx, tx, &text_value.TextValue{TextSourceID: id}); err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	if err = r.db.GormTx(tx).Delete(r.model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}

func (r *TextSourceRepository) joins(db *gorm.DB, langID uint) *gorm.DB {
	if langID == 0 {
		return db.Joins("TextValues")
	} else {
		return db.Preload("TextValue", "lang_id = ?", langID)
	}
	return db
}

func (r *TextSourceRepository) GetValuesTx(ctx context.Context, tx *gorm.DB, langID uint, sourceIDs ...*uint) ([]*string, error) {
	return r.textValueRepository.GetValuesTx(ctx, tx, langID, sourceIDs...)
}

func (r *TextSourceRepository) CreateValueTx(ctx context.Context, tx *gorm.DB, value *string, langID uint) (sourceID *uint, err error) {
	if value == nil || *value == "" || langID == 0 {
		return sourceID, nil
	}

	source := text_source.New()
	if err := r.createTx(ctx, tx, source); err != nil {
		return sourceID, err
	}
	sourceID = &source.ID

	textValue := text_value.New()
	textValue.TextSourceID = *sourceID
	textValue.TextLangID = langID
	textValue.Value = *value

	return sourceID, r.textValueRepository.CreateTx(ctx, tx, textValue)
}

func (r *TextSourceRepository) UpdateValueTx(ctx context.Context, tx *gorm.DB, sourceID *uint, value *string, langID uint) (resSourceID *uint, err error) {
	if value == nil || *value == "" || langID == 0 {
		return sourceID, nil
	}

	if sourceID == nil {
		return r.CreateValueTx(ctx, tx, value, langID)
	}

	textValue, err := r.textValueRepository.FirstTx(ctx, tx, &text_value.TextValue{
		TextSourceID: *sourceID,
		TextLangID:   langID,
	})
	textValue.Value = *value

	if err != nil {
		if errors.Is(err, yaruserror.ErrNotFound) {
			return sourceID, r.textValueRepository.CreateTx(ctx, tx, textValue)
		}
		return sourceID, err
	} else {
		if err := r.textValueRepository.UpdateTx(ctx, tx, textValue); err != nil {
			return sourceID, err
		}
	}

	return sourceID, nil
}
