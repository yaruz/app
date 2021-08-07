package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/internal/pkg/apperror"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_source"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_view_type"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"

	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// PropertyTypeRepository is a repository for the model entity
type PropertyTypeRepository struct {
	repository
	textSourceRepository text_source.Repository
}

var _ property_type.Repository = (*PropertyTypeRepository)(nil)

// New creates a new PropertyTypeRepository
func NewPropertyTypeRepository(repository *repository, textSourceRepository text_source.Repository) (*PropertyTypeRepository, error) {
	return &PropertyTypeRepository{
		repository:           *repository,
		textSourceRepository: textSourceRepository,
	}, nil
}

// Get reads the album with the specified ID from the database.
func (r *PropertyTypeRepository) Get(ctx context.Context, id uint) (*property_type.PropertyType, error) {
	return r.getTx(ctx, r.db.DB(), id)
}

func (r *PropertyTypeRepository) TGet(ctx context.Context, id uint, langID uint) (*property_type.PropertyType, error) {
	var entity *property_type.PropertyType
	err := r.db.DB().Transaction(func(tx *gorm.DB) error {
		var err error
		entity, err = r.getTx(ctx, tx, id)
		r.entityNameAndDescriptionInitTx(ctx, tx, entity, langID)
		return err
	})
	return entity, err
}

func (r *PropertyTypeRepository) entityNameAndDescriptionInitTx(ctx context.Context, tx *gorm.DB, entity *property_type.PropertyType, langID uint) error {
	s, err := r.textSourceRepository.GetValuesTx(ctx, tx, langID, entity.NameSourceID, entity.DescriptionSourceID)
	entity.Name = s[0]
	entity.Description = s[1]
	return err
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
	return r.firstTx(ctx, r.db.DB(), entity)
}

func (r *PropertyTypeRepository) TFirst(ctx context.Context, entity *property_type.PropertyType, langID uint) (*property_type.PropertyType, error) {
	err := r.db.DB().Transaction(func(tx *gorm.DB) error {
		var err error
		entity, err = r.firstTx(ctx, tx, entity)
		r.entityNameAndDescriptionInitTx(ctx, tx, entity, langID)
		return err
	})
	return entity, err
}

func (r *PropertyTypeRepository) firstTx(ctx context.Context, tx *gorm.DB, entity *property_type.PropertyType) (*property_type.PropertyType, error) {
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
func (r *PropertyTypeRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property_type.PropertyType, error) {
	return r.queryTx(ctx, r.db.DB(), cond)
}

func (r *PropertyTypeRepository) TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]property_type.PropertyType, error) {
	var items []property_type.PropertyType
	err := r.db.DB().Transaction(func(tx *gorm.DB) error {
		var err error
		items, err = r.queryTx(ctx, tx, cond)
		if err != nil {
			return err
		}

		for i := range items {
			err = r.entityNameAndDescriptionInitTx(ctx, tx, &items[i], langID)
			if err != nil {
				return err
			}
		}
		return err
	})
	return items, err
}

func (r *PropertyTypeRepository) queryTx(ctx context.Context, tx *gorm.DB, cond *selection_condition.SelectionCondition) ([]property_type.PropertyType, error) {
	items := []property_type.PropertyType{}
	db := minipkg_gorm.Conditions(tx, cond)
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
	db := minipkg_gorm.Conditions(r.db.DB(), cond)
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

func (r *PropertyTypeRepository) TCreate(ctx context.Context, entity *property_type.PropertyType, langID uint) (err error) {

	return r.db.DB().Transaction(func(tx *gorm.DB) error {
		if entity.ID > 0 {
			return errors.New("entity is not new")
		}

		if entity.NameSourceID, err = r.textSourceRepository.CreateValueTx(ctx, tx, entity.Name, langID); err != nil {
			return err
		}

		if entity.DescriptionSourceID, err = r.textSourceRepository.CreateValueTx(ctx, tx, entity.Description, langID); err != nil {
			return err
		}
		return tx.Create(entity).Error
	})
}

// Update saves a changed record in the database.
func (r *PropertyTypeRepository) Update(ctx context.Context, entity *property_type.PropertyType) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}
	return r.saveTx(ctx, r.db.DB(), entity)
}

func (r *PropertyTypeRepository) TUpdate(ctx context.Context, entity *property_type.PropertyType, langID uint) (err error) {

	return r.db.DB().Transaction(func(tx *gorm.DB) error {
		if entity.ID == 0 {
			return errors.New("entity is new")
		}

		if entity.NameSourceID, err = r.textSourceRepository.UpdateValueTx(ctx, tx, entity.NameSourceID, entity.Name, langID); err != nil {
			return err
		}

		if entity.DescriptionSourceID, err = r.textSourceRepository.UpdateValueTx(ctx, tx, entity.DescriptionSourceID, entity.Description, langID); err != nil {
			return err
		}
		return r.saveTx(ctx, tx, entity)
	})
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *PropertyTypeRepository) saveTx(ctx context.Context, tx *gorm.DB, entity *property_type.PropertyType) error {
	return tx.Save(entity).Error
}

// Delete (soft) deletes a record in the database.
func (r *PropertyTypeRepository) Delete(ctx context.Context, entity *property_type.PropertyType) error {
	return r.db.DB().Transaction(func(tx *gorm.DB) error {
		if err := r.unbindAllPropertyViewTypesTx(ctx, tx, entity); err != nil {
			return err
		}

		if err := tx.Delete(r.model, entity.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.ErrNotFound
			}
			return err
		}

		if entity.NameSourceID != nil {
			if err := r.textSourceRepository.DeleteTx(ctx, tx, *entity.NameSourceID); err != nil {
				return err
			}
		}

		if entity.DescriptionSourceID != nil {
			if err := r.textSourceRepository.DeleteTx(ctx, tx, *entity.DescriptionSourceID); err != nil {
				return err
			}
		}
		return nil
	})
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

func (r *PropertyTypeRepository) unbindAllPropertyViewTypesTx(ctx context.Context, tx *gorm.DB, entity *property_type.PropertyType) error {
	return tx.Model(entity).Association("PropertyViewTypes").Clear()
}
