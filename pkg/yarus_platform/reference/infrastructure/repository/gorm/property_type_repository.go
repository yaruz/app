package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_value"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_view_type"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"

	"github.com/yaruz/app/internal/pkg/apperror"
	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// PropertyTypeRepository is a repository for the model entity
type PropertyTypeRepository struct {
	repository
	textValueRepository text_value.Repository
}

var _ property_type.Repository = (*PropertyTypeRepository)(nil)

// New creates a new PropertyTypeRepository
func NewPropertyTypeRepository(repository *repository, textValueRepository text_value.Repository) (*PropertyTypeRepository, error) {
	return &PropertyTypeRepository{
		repository:          *repository,
		textValueRepository: textValueRepository,
	}, nil
}

// Get reads the album with the specified ID from the database.
func (r *PropertyTypeRepository) Get(ctx context.Context, id uint) (*property_type.PropertyType, error) {
	return r.getTx(ctx, r.DB(), id)
}

func (r *PropertyTypeRepository) TGet(ctx context.Context, id uint, langID uint) (*property_type.PropertyType, error) {
	var entity *property_type.PropertyType
	err := r.db.DB().Transaction(func(tx *gorm.DB) error {
		var err error
		entity, err = r.getTx(ctx, tx, id)

		if langID > 0 {
			IDs := make([]uint, 0, 2)

			if entity.NameSourceID != nil {
				IDs = append(IDs, *entity.NameSourceID)
			}

			if entity.DescriptionSourceID != nil {
				IDs = append(IDs, *entity.DescriptionSourceID)
			}
		}
		return err
	})
	return entity, err
}

func (r *PropertyTypeRepository) getTx(ctx context.Context, tx *gorm.DB, id uint) (*property_type.PropertyType, error) {
	entity := &property_type.PropertyType{}

	err := tx.First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

func (r *PropertyTypeRepository) First(ctx context.Context, entity *property_type.PropertyType) (*property_type.PropertyType, error) {
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
func (r *PropertyTypeRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property_type.PropertyType, error) {
	items := []property_type.PropertyType{}
	db := minipkg_gorm.Conditions(r.DB(), cond)
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

func (r *PropertyTypeRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
	var count int64
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.DB(), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	return count, db.Model(r.model).Count(&count).Error
}

// Create saves a new record in the database.
func (r *PropertyTypeRepository) Create(ctx context.Context, entity *property_type.PropertyType) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}

// Update saves a changed record in the database.
func (r *PropertyTypeRepository) Update(ctx context.Context, entity *property_type.PropertyType) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}
	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *PropertyTypeRepository) Save(ctx context.Context, entity *property_type.PropertyType) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a record in the database.
func (r *PropertyTypeRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(r.model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}

func (r *PropertyTypeRepository) InitPropertyViewTypes(ctx context.Context, entity *property_type.PropertyType) error {
	return r.db.DB().Model(entity).Association("PropertyViewTypes").Error
}

func (r *PropertyTypeRepository) BindPropertyViewType(ctx context.Context, entity *property_type.PropertyType, viewTypeID uint) error {
	return r.db.DB().Model(entity).Association("PropertyViewTypes").Append(&property_view_type.PropertyViewType{ID: viewTypeID})
}

func (r *PropertyTypeRepository) UnbindPropertyViewType(ctx context.Context, entity *property_type.PropertyType, viewTypeID uint) error {
	return r.db.DB().Model(entity).Association("PropertyViewTypes").Delete(&property_view_type.PropertyViewType{ID: viewTypeID})
}
