package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_value"

	"github.com/yaruz/app/internal/pkg/apperror"

	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
)

// TextValueRepository is a repository for the model entity
type TextValueRepository struct {
	repository
}

var _ text_value.Repository = (*TextValueRepository)(nil)

// New creates a new TextValueRepository
func NewTextValueRepository(repository *repository) (*TextValueRepository, error) {
	return &TextValueRepository{repository: *repository}, nil
}

// Get reads the album with the specified ID from the database.
func (r *TextValueRepository) Get(ctx context.Context, id uint) (*text_value.TextValue, error) {
	entity := &text_value.TextValue{}

	err := r.db.DB().First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruserror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

func (r *TextValueRepository) GetValuesTx(ctx context.Context, tx *gorm.DB, langID uint, sourceIDs ...*uint) ([]*string, error) {
	var err error
	values := make([]*string, len(sourceIDs))

	if langID == 0 || len(sourceIDs) == 0 {
		return values, err
	}

	IDs := make([]interface{}, 0, 2)
	for _, sourceID := range sourceIDs {
		if sourceID != nil {
			IDs = append(IDs, *sourceID)
		}
	}

	textValues, err := r.QueryTx(ctx, tx, &selection_condition.SelectionCondition{
		Where: selection_condition.WhereConditions{
			selection_condition.WhereCondition{
				Field:     "TextLangID",
				Condition: selection_condition.ConditionEq,
				Value:     langID,
			},
			selection_condition.WhereCondition{
				Field:     "TextSourceID",
				Condition: selection_condition.ConditionIn,
				Value:     IDs,
			},
		},
	})
	if err != nil {
		return values, err
	}
	valuesBySourceID := make(map[uint]*string, len(textValues))

	for i := range textValues {
		valuesBySourceID[textValues[i].TextSourceID] = &textValues[i].Value
	}

	for i, sourceID := range sourceIDs {
		if sourceID != nil {
			if val, ok := valuesBySourceID[*sourceID]; ok {
				values[i] = val
			}
		}
	}
	return values, err
}

func (r *TextValueRepository) First(ctx context.Context, entity *text_value.TextValue) (*text_value.TextValue, error) {
	return r.FirstTx(ctx, r.db.DB(), entity)
}

func (r *TextValueRepository) FirstTx(ctx context.Context, tx *gorm.DB, entity *text_value.TextValue) (*text_value.TextValue, error) {
	err := r.db.GormTx(tx).Where(entity).First(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruserror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *TextValueRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]text_value.TextValue, error) {
	return r.QueryTx(ctx, r.db.DB(), cond)
}

func (r *TextValueRepository) QueryTx(ctx context.Context, tx *gorm.DB, cond *selection_condition.SelectionCondition) ([]text_value.TextValue, error) {
	items := []text_value.TextValue{}
	db := minipkg_gorm.Conditions(r.db.GormTx(tx), cond)
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

func (r *TextValueRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
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
func (r *TextValueRepository) Create(ctx context.Context, entity *text_value.TextValue) error {
	return r.CreateTx(ctx, r.db.DB(), entity)
}

func (r *TextValueRepository) CreateTx(ctx context.Context, tx *gorm.DB, entity *text_value.TextValue) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}

	return r.db.GormTx(tx).Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *TextValueRepository) Update(ctx context.Context, entity *text_value.TextValue) error {
	return r.UpdateTx(ctx, r.db.DB(), entity)
}

func (r *TextValueRepository) UpdateTx(ctx context.Context, tx *gorm.DB, entity *text_value.TextValue) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}

	return r.saveTx(ctx, tx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *TextValueRepository) Save(ctx context.Context, entity *text_value.TextValue) error {
	return r.saveTx(ctx, r.db.DB(), entity)
}

func (r *TextValueRepository) saveTx(ctx context.Context, tx *gorm.DB, entity *text_value.TextValue) error {
	return r.db.GormTx(tx).Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *TextValueRepository) Delete(ctx context.Context, entity *text_value.TextValue) error {
	return r.DeleteTx(ctx, r.db.DB(), entity)
}

func (r *TextValueRepository) DeleteTx(ctx context.Context, tx *gorm.DB, entity *text_value.TextValue) error {

	err := r.db.GormTx(tx).Delete(r.model, entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}
