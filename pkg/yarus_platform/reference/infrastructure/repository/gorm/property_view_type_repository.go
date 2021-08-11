package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_source"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_view_type"

	"github.com/yaruz/app/internal/pkg/apperror"
	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// PropertyViewTypeRepository is a repository for the model entity
type PropertyViewTypeRepository struct {
	repository
	textSourceRepository text_source.Repository
}

var _ property_view_type.Repository = (*PropertyViewTypeRepository)(nil)

// New creates a new PropertyViewTypeRepository
func NewPropertyViewTypeRepository(repository *repository, textSourceRepository text_source.Repository) (*PropertyViewTypeRepository, error) {
	return &PropertyViewTypeRepository{
		repository:           *repository,
		textSourceRepository: textSourceRepository,
	}, nil
}

// Get reads the album with the specified ID from the database.
func (r *PropertyViewTypeRepository) Get(ctx context.Context, id uint) (*property_view_type.PropertyViewType, error) {
	return r.getTx(ctx, r.db.DB(), id)
}

func (r *PropertyViewTypeRepository) getTx(ctx context.Context, tx *gorm.DB, id uint) (*property_view_type.PropertyViewType, error) {
	entity := &property_view_type.PropertyViewType{}

	err := tx.First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

func (r *PropertyViewTypeRepository) First(ctx context.Context, entity *property_view_type.PropertyViewType) (*property_view_type.PropertyViewType, error) {
	return r.firstTx(ctx, r.db.DB(), entity)
}

func (r *PropertyViewTypeRepository) firstTx(ctx context.Context, tx *gorm.DB, entity *property_view_type.PropertyViewType) (*property_view_type.PropertyViewType, error) {
	err := tx.Where(entity).First(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *PropertyViewTypeRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property_view_type.PropertyViewType, error) {
	items := []property_view_type.PropertyViewType{}
	db := minipkg_gorm.Conditions(r.db.DB(), cond)
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

func (r *PropertyViewTypeRepository) entityNameAndDescriptionInitTx(ctx context.Context, tx *gorm.DB, entity *property_view_type.PropertyViewType, langID uint) error {
	s, err := r.textSourceRepository.GetValuesTx(ctx, tx, langID, entity.NameSourceID, entity.DescriptionSourceID)
	entity.Name = s[0]
	entity.Description = s[1]
	return err
}

func (r *PropertyViewTypeRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
	var count int64
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.db.DB(), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	err := db.Count(&count).Error
	return count, err
}

// Create saves a new record in the database.
func (r *PropertyViewTypeRepository) Create(ctx context.Context, entity *property_view_type.PropertyViewType) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *PropertyViewTypeRepository) Update(ctx context.Context, entity *property_view_type.PropertyViewType) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}
	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *PropertyViewTypeRepository) Save(ctx context.Context, entity *property_view_type.PropertyViewType) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *PropertyViewTypeRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(r.model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}
