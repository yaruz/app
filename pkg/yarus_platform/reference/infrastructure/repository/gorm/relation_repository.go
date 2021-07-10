package gorm

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"

	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property_type"

	"github.com/pkg/errors"

	minipkg_gorm "github.com/minipkg/db/gorm"
	"github.com/minipkg/selection_condition"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type2property"
	"github.com/yaruz/app/pkg/yarus_platform/yaruzerror"
	"gorm.io/gorm"
)

// RelationRepository is a repository for the model entity
type RelationRepository struct {
	repository
	entityType2PropertyRepository entity_type2property.Repository
}

var _ entity_type.RelationRepository = (*RelationRepository)(nil)

// New creates a new RelationRepository
func NewRelationRepository(repository *repository, entityType2PropertyRepository *entity_type2property.Repository) (*RelationRepository, error) {
	return &RelationRepository{repository: *repository, entityType2PropertyRepository: *entityType2PropertyRepository}, nil
}

// Get reads the album with the specified ID from the database.
func (r *RelationRepository) Get(ctx context.Context, id uint) (*entity_type.Relation, error) {
	entity := &entity_type.Relation{}

	err := r.joins(r.relationTypeDB()).First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}

	if err = r.AfterFind(ctx, entity); err != nil {
		return nil, err
	}

	return entity, err
}

func (r *RelationRepository) First(ctx context.Context, entity *entity_type.Relation) (*entity_type.Relation, error) {
	err := r.joins(r.relationTypeDB()).Where(entity).First(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, yaruzerror.ErrNotFound
		}
		return nil, err
	}

	if err = r.AfterFind(ctx, entity); err != nil {
		return nil, err
	}

	return entity, err
}

// Query retrieves the album records with the specified offset and limit from the database.
func (r *RelationRepository) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]entity_type.Relation, error) {
	items := []entity_type.Relation{}
	db := minipkg_gorm.Conditions(r.relationTypeDB(), cond)
	if db.Error != nil {
		return nil, db.Error
	}

	err := r.joins(db).Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, yaruzerror.ErrNotFound
		}
		return nil, err
	}

	for _, entity := range items {
		if err = r.AfterFind(ctx, &entity); err != nil {
			return nil, err
		}
	}

	return items, err
}

func (r *RelationRepository) PropertyAndRelationQuery(ctx context.Context, cond *selection_condition.SelectionCondition) ([]property.Property, []entity_type.Relation, error) {
	items := []property.Property{}
	db := minipkg_gorm.Conditions(r.DB(), cond)
	if db.Error != nil {
		return nil, nil, db.Error
	}

	err := r.joins(db).Find(&items).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, yaruzerror.ErrNotFound
		}
		return nil, nil, err
	}
	props := make([]property.Property, 0, len(items))
	rels := make([]entity_type.Relation, 0, len(items))

	for _, item := range items {

		if item.PropertyTypeID == property_type.IDRelation {
			rel := entity_type.NewRelation()
			rel.Property = item

			if err = r.AfterFind(ctx, rel); err != nil {
				return nil, nil, err
			}
			rels = append(rels, *rel)
		} else {

			if err = item.AfterFind(); err != nil {
				return nil, nil, err
			}
			props = append(props, item)
		}
	}

	return props, rels, err
}

func (r *RelationRepository) GetPropertiesAndRelationsByEntityTypeID(ctx context.Context, entityTypeID uint) ([]property.Property, []entity_type.Relation, error) {
	rels, err := r.entityType2PropertyRepository.Query(ctx, &selection_condition.SelectionCondition{
		Where: &entity_type2property.EntityType2Property{
			EntityTypeID: entityTypeID,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	IDs := make([]interface{}, 0, len(rels))
	for _, rel := range rels {
		IDs = append(IDs, rel.PropertyID)
	}

	return r.PropertyAndRelationQuery(ctx, &selection_condition.SelectionCondition{
		Where: selection_condition.WhereCondition{
			Field:     "ID",
			Condition: selection_condition.ConditionIn,
			Value:     IDs,
		},
	})
}

func (r *RelationRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
	var count int64
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.relationTypeDB(), cond)
	if db.Error != nil {
		return 0, db.Error
	}

	err := db.Count(&count).Error
	return count, err
}

func (r *RelationRepository) AfterFind(ctx context.Context, entity *entity_type.Relation) error {

	if err := entity.AfterFind(); err != nil {
		return err
	}

	if err := r.InitRelatedEntityTypes(ctx, entity); err != nil {
		return err
	}
	return nil
}

func (r *RelationRepository) InitRelatedEntityTypes(ctx context.Context, entity *entity_type.Relation) error {
	rels, err := r.entityType2PropertyRepository.QueryWithEntityType(ctx, &selection_condition.SelectionCondition{
		Where: &entity_type2property.EntityType2Property{
			PropertyID: entity.ID,
		},
	})
	if err != nil {
		return err
	}
	for _, r := range rels {
		if r.IsDependent {
			entity.DependedEntityType = r.EntityType
		} else {
			entity.UndependedEntityType = r.EntityType
		}
	}
	return nil
}

// Create saves a new record in the database.
func (r *RelationRepository) Create(ctx context.Context, entity *entity_type.Relation) error {

	if entity.ID > 0 {
		return errors.New("entity is not new")
	}

	if err := entity.BeforeSave(); err != nil {
		return err
	}

	return r.db.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(entity).Error; err != nil {
			return err
		}

		return r.createLinksTx(ctx, tx, entity)
	})
}

// Update saves a changed Maintenance record in the database.
func (r *RelationRepository) Update(ctx context.Context, entity *entity_type.Relation) error {

	if entity.ID == 0 {
		return errors.New("entity is new")
	}

	return r.Save(ctx, entity)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (r *RelationRepository) Save(ctx context.Context, entity *entity_type.Relation) error {

	if entity.ID == 0 {
		return r.Create(ctx, entity)
	}

	if err := entity.BeforeSave(); err != nil {
		return err
	}

	return r.db.DB().Transaction(func(tx *gorm.DB) error {
		if err := r.deleteAllLinksTx(ctx, tx, entity.ID); err != nil {
			return err
		}

		if err := r.createLinksTx(ctx, tx, entity); err != nil {
			return err
		}

		return tx.Save(entity).Error
	})
}

// Delete (soft) deletes a Maintenance record in the database.
func (r *RelationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.DB().Transaction(func(tx *gorm.DB) error {
		if err := r.deleteAllLinksTx(ctx, tx, id); err != nil {
			return err
		}

		err := tx.Delete(r.model, id).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.ErrNotFound
			}
		}

		return err
	})
}

func (r *RelationRepository) deleteAllLinksTx(ctx context.Context, tx *gorm.DB, id uint) error {
	return r.entityType2PropertyRepository.DeleteTx(ctx, tx, &entity_type2property.EntityType2Property{
		PropertyID: id,
	})
}

func (r *RelationRepository) createLinksTx(ctx context.Context, tx *gorm.DB, entity *entity_type.Relation) error {
	if err := r.validateLinks(entity); err != nil {
		return err
	}

	if err := r.entityType2PropertyRepository.CreateTx(ctx, tx, &entity_type2property.EntityType2Property{
		EntityTypeID: entity.UndependedEntityType.ID,
		PropertyID:   entity.ID,
		IsDependent:  false,
	}); err != nil {
		return err
	}

	return r.entityType2PropertyRepository.CreateTx(ctx, tx, &entity_type2property.EntityType2Property{
		EntityTypeID: entity.DependedEntityType.ID,
		PropertyID:   entity.ID,
		IsDependent:  true,
	})
}

func (r *RelationRepository) validateLinks(entity *entity_type.Relation) error {
	if entity.DependedEntityType.ID == 0 || entity.UndependedEntityType.ID == 0 {
		return errors.Wrapf(apperror.ErrBadParams, "Some of links is empty")
	}
	return nil
}

func (r *RelationRepository) relationTypeDB() *gorm.DB {
	return r.DB().Where(&property.Property{
		PropertyTypeID: property_type.IDRelation,
	})
}

func (r *RelationRepository) joins(db *gorm.DB) *gorm.DB {
	return db.Joins("PropertyType").Joins("PropertyViewType").Joins("PropertyGroup")
}
