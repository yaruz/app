package gorm

import (
	"context"

	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_source"

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
	propertyRepository            property.Repository
	entityType2PropertyRepository entity_type2property.Repository
	textSourceRepository          text_source.Repository
}

var _ entity_type.RelationRepository = (*RelationRepository)(nil)

// New creates a new RelationRepository
func NewRelationRepository(repository *repository, propertyRepository property.Repository, entityType2PropertyRepository entity_type2property.Repository, textSourceRepository text_source.Repository) (*RelationRepository, error) {
	return &RelationRepository{
		repository:                    *repository,
		propertyRepository:            propertyRepository,
		entityType2PropertyRepository: entityType2PropertyRepository,
		textSourceRepository:          textSourceRepository,
	}, nil
}

// Get reads the album with the specified ID from the database.
func (r *RelationRepository) Get(ctx context.Context, id uint) (*entity_type.Relation, error) {
	return r.getTx(ctx, r.db.DB(), id)
}

func (r *RelationRepository) TGet(ctx context.Context, id uint, langID uint) (*entity_type.Relation, error) {
	var entity *entity_type.Relation
	err := r.db.DB().Transaction(func(tx *gorm.DB) error {
		var err error
		entity, err = r.getTx(ctx, tx, id)
		r.entityNameAndDescriptionInitTx(ctx, tx, entity, langID)
		return err
	})
	return entity, err
}

func (r *RelationRepository) getTx(ctx context.Context, tx *gorm.DB, id uint) (*entity_type.Relation, error) {
	entity := &entity_type.Relation{}

	err := r.joins(r.propertyTypeRelationTx(tx)).First(entity, id).Error
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
	return r.firstTx(ctx, r.db.DB(), entity)
}

func (r *RelationRepository) TFirst(ctx context.Context, entity *entity_type.Relation, langID uint) (*entity_type.Relation, error) {
	err := r.db.DB().Transaction(func(tx *gorm.DB) error {
		var err error
		entity, err = r.firstTx(ctx, tx, entity)
		r.entityNameAndDescriptionInitTx(ctx, tx, entity, langID)
		return err
	})
	return entity, err
}

func (r *RelationRepository) firstTx(ctx context.Context, tx *gorm.DB, entity *entity_type.Relation) (*entity_type.Relation, error) {
	err := r.joins(r.propertyTypeRelationTx(tx)).Where(entity).First(entity).Error
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
	return r.queryTx(ctx, r.db.DB(), cond)
}

func (r *RelationRepository) TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]entity_type.Relation, error) {
	var items []entity_type.Relation
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

func (r *RelationRepository) queryTx(ctx context.Context, tx *gorm.DB, cond *selection_condition.SelectionCondition) ([]entity_type.Relation, error) {
	items := []entity_type.Relation{}
	db := minipkg_gorm.Conditions(r.propertyTypeRelationTx(tx), cond)
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
	return r.propertyAndRelationQueryTx(ctx, r.db.DB(), cond)
}

func (r *RelationRepository) TPropertyAndRelationQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]property.Property, []entity_type.Relation, error) {
	return r.tPropertyAndRelationQueryTx(ctx, r.db.DB(), cond, langID)
}

func (r *RelationRepository) tPropertyAndRelationQueryTx(ctx context.Context, tx *gorm.DB, cond *selection_condition.SelectionCondition, langID uint) ([]property.Property, []entity_type.Relation, error) {
	var rels []entity_type.Relation
	var props []property.Property
	err := tx.Transaction(func(tx *gorm.DB) error {
		var err error
		props, rels, err = r.propertyAndRelationQueryTx(ctx, tx, cond)
		if err != nil {
			return err
		}

		for i := range rels {
			err = r.entityNameAndDescriptionInitTx(ctx, tx, &rels[i], langID)
			if err != nil {
				return err
			}
		}

		for p := range props {
			err = r.propertyRepository.EntityNameAndDescriptionInitTx(ctx, tx, &props[p], langID)
			if err != nil {
				return err
			}
		}
		return err
	})
	return props, rels, err
}

func (r *RelationRepository) propertyAndRelationQueryTx(ctx context.Context, tx *gorm.DB, cond *selection_condition.SelectionCondition) ([]property.Property, []entity_type.Relation, error) {
	items := []property.Property{}
	db := minipkg_gorm.Conditions(tx, cond)
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

func (r *RelationRepository) getPropertiesIDsByEntityTypeIDTx(ctx context.Context, tx *gorm.DB, entityTypeID uint) ([]interface{}, error) {
	rs, err := r.entityType2PropertyRepository.QueryTx(ctx, tx, &selection_condition.SelectionCondition{
		Where: &entity_type2property.EntityType2Property{
			EntityTypeID: entityTypeID,
		},
	})
	if err != nil {
		return nil, err
	}

	IDs := make([]interface{}, 0, len(rs))
	for _, rel := range rs {
		IDs = append(IDs, rel.PropertyID)
	}
	return IDs, nil
}

func (r *RelationRepository) GetPropertiesAndRelationsByEntityTypeID(ctx context.Context, entityTypeID uint) ([]property.Property, []entity_type.Relation, error) {

	var rels []entity_type.Relation
	var props []property.Property
	err := r.db.DB().Transaction(func(tx *gorm.DB) error {
		IDs, err := r.getPropertiesIDsByEntityTypeIDTx(ctx, tx, entityTypeID)
		if err != nil {
			return err
		}

		props, rels, err = r.propertyAndRelationQueryTx(ctx, tx, &selection_condition.SelectionCondition{
			Where: selection_condition.WhereCondition{
				Field:     "ID",
				Condition: selection_condition.ConditionIn,
				Value:     IDs,
			},
		})
		return err
	})
	return props, rels, err
}

func (r *RelationRepository) TGetPropertiesAndRelationsByEntityTypeID(ctx context.Context, entityTypeID uint, langID uint) ([]property.Property, []entity_type.Relation, error) {

	var rels []entity_type.Relation
	var props []property.Property
	err := r.db.DB().Transaction(func(tx *gorm.DB) error {
		IDs, err := r.getPropertiesIDsByEntityTypeIDTx(ctx, tx, entityTypeID)
		if err != nil {
			return err
		}

		props, rels, err = r.tPropertyAndRelationQueryTx(ctx, tx, &selection_condition.SelectionCondition{
			Where: selection_condition.WhereCondition{
				Field:     "ID",
				Condition: selection_condition.ConditionIn,
				Value:     IDs,
			},
		}, langID)
		return err
	})
	return props, rels, err
}

func (r *RelationRepository) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
	var count int64
	c := cond
	c.Limit = 0
	c.Offset = 0
	db := minipkg_gorm.Conditions(r.propertyTypeRelationTx(r.db.DB()), cond)
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

func (r *RelationRepository) entityNameAndDescriptionInitTx(ctx context.Context, tx *gorm.DB, entity *entity_type.Relation, langID uint) error {
	s, err := r.textSourceRepository.GetValuesTx(ctx, tx, langID, entity.NameSourceID, entity.DescriptionSourceID)
	entity.Name = s[0]
	entity.Description = s[1]
	return err
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

func (r *RelationRepository) TCreate(ctx context.Context, entity *entity_type.Relation, langID uint) (err error) {

	if entity.ID > 0 {
		return errors.New("entity is not new")
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

	return r.saveTx(ctx, r.db.DB(), entity)
}

func (r *RelationRepository) TUpdate(ctx context.Context, entity *entity_type.Relation, langID uint) (err error) {

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
func (r *RelationRepository) saveTx(ctx context.Context, tx *gorm.DB, entity *entity_type.Relation) error {

	if entity.ID == 0 {
		return r.Create(ctx, entity)
	}

	if err := entity.BeforeSave(); err != nil {
		return err
	}

	return tx.Transaction(func(tx *gorm.DB) error {
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
func (r *RelationRepository) Delete(ctx context.Context, entity *entity_type.Relation) error {
	return r.db.DB().Transaction(func(tx *gorm.DB) error {
		if err := r.deleteAllLinksTx(ctx, tx, entity.ID); err != nil {
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

func (r *RelationRepository) propertyTypeRelationTx(tx *gorm.DB) *gorm.DB {
	return tx.Where(&property.Property{
		PropertyTypeID: property_type.IDRelation,
	})
}

func (r *RelationRepository) joins(db *gorm.DB) *gorm.DB {
	return db.Joins("PropertyType").Joins("PropertyViewType").Joins("PropertyGroup")
}
