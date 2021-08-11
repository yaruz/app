package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_source"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_group"

	"github.com/yaruz/app/internal/pkg/apperror"
	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// PropertyGroupRepository is a repository for the model entity
type PropertyGroupRepository struct {
	repository
	textSourceRepository text_source.Repository
}

var _ property_group.Repository = (*PropertyGroupRepository)(nil)

// New creates a new PropertyGroupRepository
func NewPropertyGroupRepository(repository *repository, textSourceRepository text_source.Repository) (*PropertyGroupRepository, error) {
	return &PropertyGroupRepository{
		repository:           *repository,
		textSourceRepository: textSourceRepository,
	}, nil
}

// Get reads the album with the specified ID from the database.
func (r *PropertyGroupRepository) Get(ctx context.Context, id uint) (*property_group.PropertyGroup, error) {
	return r.getTx(ctx, r.db.DB(), id)
}

func (r *PropertyGroupRepository) getTx(ctx context.Context, tx *gorm.DB, id uint) (*property_group.PropertyGroup, error) {
	entity := &property_group.PropertyGroup{}

	err := tx.First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}
	return entity, err
}

func (r *PropertyGroupRepository) First(ctx context.Context, entity *property_group.PropertyGroup) (*property_group.PropertyGroup, error) {
	return r.firstTx(ctx, r.db.DB(), entity)
}

func (r *PropertyGroupRepository) firstTx(ctx context.Context, tx *gorm.DB, entity *property_group.PropertyGroup) (*property_group.PropertyGroup, error) {
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
func (r *PropertyGroupRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property_group.PropertyGroup, error) {
	items := []property_group.PropertyGroup{}
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

func (r *PropertyGroupRepository) entityNameAndDescriptionInitTx(ctx context.Context, tx *gorm.DB, entity *property_group.PropertyGroup, langID uint) error {
	s, err := r.textSourceRepository.GetValuesTx(ctx, tx, langID, entity.NameSourceID, entity.DescriptionSourceID)
	entity.Name = s[0]
	entity.Description = s[1]
	return err
}

func (r *PropertyGroupRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
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
func (r *PropertyGroupRepository) Create(ctx context.Context, entity *property_group.PropertyGroup) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}
	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *PropertyGroupRepository) Update(ctx context.Context, entity *property_group.PropertyGroup) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}
	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *PropertyGroupRepository) Save(ctx context.Context, entity *property_group.PropertyGroup) error {
	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *PropertyGroupRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(r.model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}
