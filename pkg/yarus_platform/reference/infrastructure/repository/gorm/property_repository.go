package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_source"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"

	"github.com/yaruz/app/internal/pkg/apperror"
	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
)

// PropertyRepository is a repository for the model entity
type PropertyRepository struct {
	repository
	textSourceRepository text_source.Repository
}

var _ property.Repository = (*PropertyRepository)(nil)

// New creates a new PropertyRepository
func NewPropertyRepository(repository *repository, textSourceRepository text_source.Repository) (*PropertyRepository, error) {
	return &PropertyRepository{
		repository:           *repository,
		textSourceRepository: textSourceRepository,
	}, nil
}

// Get reads the album with the specified ID from the database.
func (r *PropertyRepository) Get(ctx context.Context, id uint) (*property.Property, error) {
	return r.getTx(ctx, r.db.DB(), id)
}

func (r *PropertyRepository) TGet(ctx context.Context, id uint, langID uint) (*property.Property, error) {
	var entity *property.Property
	err := r.db.DB().Transaction(func(tx *gorm.DB) error {
		var err error
		entity, err = r.getTx(ctx, tx, id)
		r.EntityNameAndDescriptionInitTx(ctx, tx, entity, langID)
		return err
	})
	return entity, err
}

func (r *PropertyRepository) getTx(ctx context.Context, tx *gorm.DB, id uint) (*property.Property, error) {
	entity := &property.Property{}

	err := r.joins(r.db.GormTx(tx)).First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruserror.ErrNotFound
		}
		return nil, err
	}

	if err = entity.AfterFind(); err != nil {
		return nil, err
	}

	return entity, err
}

func (r *PropertyRepository) First(ctx context.Context, entity *property.Property) (*property.Property, error) {
	return r.firstTx(ctx, r.db.DB(), entity)
}

func (r *PropertyRepository) TFirst(ctx context.Context, entity *property.Property, langID uint) (*property.Property, error) {
	err := r.db.DB().Transaction(func(tx *gorm.DB) error {
		var err error
		entity, err = r.firstTx(ctx, tx, entity)
		r.EntityNameAndDescriptionInitTx(ctx, tx, entity, langID)
		return err
	})
	return entity, err
}

func (r *PropertyRepository) firstTx(ctx context.Context, tx *gorm.DB, entity *property.Property) (*property.Property, error) {
	err := r.joins(r.db.GormTx(tx)).Where(entity).First(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruserror.ErrNotFound
		}
		return nil, err
	}

	if err = entity.AfterFind(); err != nil {
		return nil, err
	}

	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *PropertyRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property.Property, error) {
	return r.queryTx(ctx, r.db.DB(), cond)
}

func (r *PropertyRepository) TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]property.Property, error) {
	var items []property.Property
	err := r.db.DB().Transaction(func(tx *gorm.DB) error {
		var err error
		items, err = r.queryTx(ctx, tx, cond)
		if err != nil {
			return err
		}

		for i := range items {
			err = r.EntityNameAndDescriptionInitTx(ctx, tx, &items[i], langID)
			if err != nil {
				return err
			}
		}
		return err
	})
	return items, err
}

func (r *PropertyRepository) queryTx(ctx context.Context, tx *gorm.DB, cond *selection_condition.SelectionCondition) ([]property.Property, error) {
	items := []property.Property{}
	db := minipkg_gorm.Conditions(r.db.GormTx(tx), cond)
	if db.Error != nil {
		return nil, db.Error
	}

	err := r.joins(db).Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return items, yaruserror.ErrNotFound
		}
		return nil, err
	}

	for _, entity := range items {
		if err = entity.AfterFind(); err != nil {
			return nil, err
		}
	}

	return items, err
}

func (r *PropertyRepository) EntityNameAndDescriptionInitTx(ctx context.Context, tx *gorm.DB, entity *property.Property, langID uint) error {
	s, err := r.textSourceRepository.GetValuesTx(ctx, tx, langID, entity.NameSourceID, entity.DescriptionSourceID)
	entity.Name = s[0]
	entity.Description = s[1]
	return err
}

func (r *PropertyRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
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
func (r *PropertyRepository) Create(ctx context.Context, entity *property.Property) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}

	if err := entity.Validate(); err != nil {
		return err
	}

	if err := entity.BeforeSave(); err != nil {
		return err
	}

	return r.db.DB().Create(entity).Error
}

func (r *PropertyRepository) TCreate(ctx context.Context, entity *property.Property, langID uint) (err error) {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}

	if err := entity.Validate(); err != nil {
		return err
	}

	if err := entity.BeforeSave(); err != nil {
		return err
	}

	return r.db.DB().Transaction(func(tx *gorm.DB) error {

		if entity.NameSourceID, err = r.textSourceRepository.CreateValueTx(ctx, tx, entity.Name, langID); err != nil {
			return err
		}

		if entity.DescriptionSourceID, err = r.textSourceRepository.CreateValueTx(ctx, tx, entity.Description, langID); err != nil {
			return err
		}
		return tx.Create(entity).Error
	})
}

// Update saves a changed Maintenance record in the database.
func (r *PropertyRepository) Update(ctx context.Context, entity *property.Property) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}

	if err := entity.Validate(); err != nil {
		return err
	}

	return r.saveTx(ctx, r.db.DB(), entity)
}

func (r *PropertyRepository) TUpdate(ctx context.Context, entity *property.Property, langID uint) (err error) {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}

	if err := entity.Validate(); err != nil {
		return err
	}

	return r.db.DB().Transaction(func(tx *gorm.DB) error {

		if entity.NameSourceID, err = r.textSourceRepository.UpdateValueTx(ctx, tx, entity.NameSourceID, entity.Name, langID); err != nil {
			return err
		}

		if entity.DescriptionSourceID, err = r.textSourceRepository.UpdateValueTx(ctx, tx, entity.DescriptionSourceID, entity.Description, langID); err != nil {
			return err
		}
		return r.saveTx(ctx, tx, entity)
	})
}

// saveTx update value in database, if the value doesn't have primary key, will insert it
func (r *PropertyRepository) saveTx(ctx context.Context, tx *gorm.DB, entity *property.Property) error {

	if err := entity.BeforeSave(); err != nil {
		return err
	}

	return r.db.GormTx(tx).Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *PropertyRepository) Delete(ctx context.Context, entity *property.Property) error {

	return r.db.DB().Transaction(func(tx *gorm.DB) error {

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

func (r *PropertyRepository) joins(db *gorm.DB) *gorm.DB {
	return db.Joins("PropertyType").Joins("PropertyViewType").Joins("PropertyUnit").Joins("PropertyGroup")
}
