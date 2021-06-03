package gorm

import (
	"context"
	"errors"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type2property"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/relation"

	"github.com/yaruz/app/internal/pkg/apperror"
	"gorm.io/gorm"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
)

// RelationRepository is a repository for the model entity
type RelationRepository struct {
	repository
	entityType2PropertyRepository entity_type2property.Repository
}

var _ relation.Repository = (*RelationRepository)(nil)

// New creates a new RelationRepository
func NewRelationRepository(repository *repository, entityType2PropertyRepository entity_type2property.Repository) (*RelationRepository, error) {
	return &RelationRepository{repository: *repository, entityType2PropertyRepository: entityType2PropertyRepository}, nil
}

// Get reads the album with the specified ID from the database.
func (r *RelationRepository) Get(ctx context.Context, id uint) (*relation.Relation, error) {
	entity := &relation.Relation{}

	err := r.joins(r.DB()).First(entity, id).Error
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

func (r *RelationRepository) First(ctx context.Context, entity *relation.Relation) (*relation.Relation, error) {
	err := r.joins(r.DB()).Where(entity).First(entity).Error
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
func (r *RelationRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]relation.Relation, error) {
	items := []relation.Relation{}
	db := minipkg_gorm.Conditions(r.DB(), cond)
	if db.Error != nil {
		return nil, db.Error
	}

	err := r.joins(db).Find(&items).Error
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

func (r *RelationRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
	var count int64
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.DB(), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	err := db.Count(&count).Error
	return count, err
}

// Create saves a new record in the database.
func (r *RelationRepository) Create(ctx context.Context, entity *relation.Relation) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}

	if err := entity.BeforeSave(); err != nil {
		return err
	}

	//return r.db.DB().Create(entity).Error
	return r.db.DB().Transaction(func(tx *gorm.DB) error {
		if err := r.entityType2PropertyRepository.CreateTx(ctx, tx, &entity_type2property.EntityType2Property{
			EntityTypeID: entity.UndependedEntityType.ID,
			PropertyID:   entity.ID,
			IsDependent:  false,
		}); err != nil {
			return err
		}

		if err := r.entityType2PropertyRepository.CreateTx(ctx, tx, &entity_type2property.EntityType2Property{
			EntityTypeID: entity.DependedEntityType.ID,
			PropertyID:   entity.ID,
			IsDependent:  true,
		}); err != nil {
			return err
		}

		return tx.Create(entity).Error
	})
}

// Update saves a changed Maintenance record in the database.
func (r *RelationRepository) Update(ctx context.Context, entity *relation.Relation) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}

	if err := entity.BeforeSave(); err != nil {
		return err
	}

	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *RelationRepository) Save(ctx context.Context, entity *relation.Relation) error {

	if err := entity.BeforeSave(); err != nil {
		return err
	}

	return r.db.DB().Save(entity).Error
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *RelationRepository) Delete(ctx context.Context, id uint) error {

	err := r.db.DB().Delete(r.model, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrNotFound
		}
	}
	return err
}

func (r *RelationRepository) joins(db *gorm.DB) *gorm.DB {
	return db.Joins("PropertyType").Joins("PropertyViewType").Joins("PropertyUnit").Joins("PropertyGroup")
}
