package gorm

import (
	"context"
	"errors"

	domain_entity "github.com/yaruz/app/pkg/yarus_platform/data/domain/entity"

	"github.com/yaruz/app/internal/pkg/apperror"

	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// EntityRepository is a repository for the model entity
type EntityRepository struct {
	repository
}

var _ domain_entity.Repository = (*EntityRepository)(nil)

// New creates a new EntityRepository
func NewEntityRepository(repository *repository) (*EntityRepository, error) {
	r := &EntityRepository{repository: *repository}
	r.autoMigrate()
	return r, nil
}

func (r *EntityRepository) autoMigrate() {
	if r.db.IsAutoMigrate() {
		r.db.DB().AutoMigrate(&domain_entity.Entity{})
	}
}

// Get reads the album with the specified ID from the database.
func (r *EntityRepository) Get(ctx context.Context, id uint) (*domain_entity.Entity, error) {
	entity := &domain_entity.Entity{}

	err := r.DB().First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}

	if err = entity.AfterFind(); err != nil {
		return nil, err
	}

	return entity, err
}

func (r *EntityRepository) First(ctx context.Context, entity *domain_entity.Entity) (*domain_entity.Entity, error) {
	err := r.DB().Where(entity).First(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}

	if err = entity.AfterFind(); err != nil {
		return nil, err
	}

	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *EntityRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]domain_entity.Entity, error) {
	items := []domain_entity.Entity{}
	db := minipkg_gorm.Conditions(r.DB().Model(&domain_entity.Entity{}), cond)
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

	for _, entity := range items {
		if err = entity.AfterFind(); err != nil {
			return nil, err
		}
	}

	return items, err
}

func (r *EntityRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
	var count int64
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.DB().Model(&domain_entity.Entity{}), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	err := db.Count(&count).Error
	return count, err
}

// Create saves a new record in the database.
func (r *EntityRepository) Create(ctx context.Context, entity *domain_entity.Entity) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}

	if err := entity.BeforeSave(); err != nil {
		return err
	}

	return r.db.DB().Create(entity).Error
}

// Update saves a changed Maintenance record in the database.
func (r *EntityRepository) Update(ctx context.Context, entity *domain_entity.Entity) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}

	if err := entity.BeforeSave(); err != nil {
		return err
	}

	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *EntityRepository) Save(ctx context.Context, entity *domain_entity.Entity) error {

	if err := entity.BeforeSave(); err != nil {
		return err
	}

	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *EntityRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(&domain_entity.Entity{}, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}
