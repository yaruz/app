package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_source"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_unit"

	"github.com/yaruz/app/internal/pkg/apperror"
	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// PropertyUnitRepository is a repository for the model entity
type PropertyUnitRepository struct {
	repository
	textSourceRepository text_source.Repository
}

var _ property_unit.Repository = (*PropertyUnitRepository)(nil)

// New creates a new PropertyUnitRepository
func NewPropertyUnitRepository(repository *repository, textSourceRepository text_source.Repository) (*PropertyUnitRepository, error) {
	return &PropertyUnitRepository{
		repository:           *repository,
		textSourceRepository: textSourceRepository,
	}, nil
}

// Get reads the album with the specified ID from the database.
func (r *PropertyUnitRepository) Get(ctx context.Context, id uint) (*property_unit.PropertyUnit, error) {
	return r.getTx(ctx, r.db.DB(), id)
}

func (r *PropertyUnitRepository) getTx(ctx context.Context, tx *gorm.DB, id uint) (*property_unit.PropertyUnit, error) {
	entity := &property_unit.PropertyUnit{}

	err := tx.First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

func (r *PropertyUnitRepository) First(ctx context.Context, entity *property_unit.PropertyUnit) (*property_unit.PropertyUnit, error) {
	return r.firstTx(ctx, r.db.DB(), entity)
}

func (r *PropertyUnitRepository) firstTx(ctx context.Context, tx *gorm.DB, entity *property_unit.PropertyUnit) (*property_unit.PropertyUnit, error) {
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
func (r *PropertyUnitRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property_unit.PropertyUnit, error) {
	items := []property_unit.PropertyUnit{}
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

func (r *PropertyUnitRepository) entityNameAndDescriptionInitTx(ctx context.Context, tx *gorm.DB, entity *property_unit.PropertyUnit, langID uint) error {
	s, err := r.textSourceRepository.GetValuesTx(ctx, tx, langID, entity.NameSourceID, entity.DescriptionSourceID)
	entity.Name = s[0]
	entity.Description = s[1]
	return err
}

func (r *PropertyUnitRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
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
func (r *PropertyUnitRepository) Create(ctx context.Context, entity *property_unit.PropertyUnit) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *PropertyUnitRepository) Update(ctx context.Context, entity *property_unit.PropertyUnit) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}
	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *PropertyUnitRepository) Save(ctx context.Context, entity *property_unit.PropertyUnit) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *PropertyUnitRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(r.model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}
