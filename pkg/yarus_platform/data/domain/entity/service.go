package entity

import (
	"context"

	"github.com/pkg/errors"

	"github.com/minipkg/log"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/reference"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"
	"github.com/yaruz/app/pkg/yarus_platform/search"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
)

// IService encapsulates usecase logic.
type IService interface {
	NewEntity() *Entity
	Get(ctx context.Context, id uint, langID uint) (*Entity, error)
	First(ctx context.Context, entity *Entity, langID uint) (*Entity, error)
	Query(ctx context.Context, query *selection_condition.SelectionCondition, langID uint) ([]Entity, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *Entity, langID uint) error
	Update(ctx context.Context, entity *Entity, langID uint) error
	Delete(ctx context.Context, id uint) error
	EntityInit(ctx context.Context, entity *Entity, langID uint) error
	//SetValueForProperty(entity *Entity, property *property.Property, value interface{}, langID uint) error
	//SetValueForRelation(entity *Entity, relation *entity_type.Relation, value []uint) error
	//EntitySetPropertyValue(entity *Entity, propertyValue *PropertyValue)
	//EntitySetRelationValue(entity *Entity, relationValue *RelationValue)
	//EntityDeletePropertyValue(entity *Entity, propertyID uint)
	BindRelatedEntity(relation *entity_type.Relation, entity1 *Entity, entity2 *Entity) error
	BindRelatedEntities(relation *entity_type.Relation, entity *Entity, relatedEntities []*Entity) error
	UnbindRelatedEntity(relation *entity_type.Relation, entity1 *Entity, entity2 *Entity) error
	UnbindRelatedEntities(relation *entity_type.Relation, entity *Entity, relatedEntities []*Entity) error
	CheckBindRelatedEntity(relation *entity_type.Relation, entity1 *Entity, entity2 *Entity) (firstParamEntityIsDepended bool, err error)
	CheckBindRelatedEntities(relation *entity_type.Relation, entity *Entity, relatedEntities []*Entity) (firstParamEntityIsDepended bool, err error)
}

type service struct {
	logger     log.ILogger
	repository Repository
	reference  *reference.ReferenceSubsystem
	search     *search.SearchSubsystem
}

var _ IService = (*service)(nil)

// NewService creates a new service.
func NewService(logger log.ILogger, repo Repository, reference *reference.ReferenceSubsystem, search *search.SearchSubsystem) IService {
	s := &service{
		logger:     logger,
		repository: repo,
		reference:  reference,
		search:     search,
	}
	//repo.SetDefaultConditions(s.defaultConditions())
	return s
}

func (s *service) NewEntity() *Entity {
	return New()
}

// Get returns the entity with the specified ID.
func (s *service) Get(ctx context.Context, id uint, langID uint) (*Entity, error) {
	entity, err := s.repository.Get(ctx, id, langID)
	if err != nil {
		return nil, err
	}

	if err = s.EntityInit(ctx, entity, langID); err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *service) First(ctx context.Context, entity *Entity, langID uint) (*Entity, error) {
	entity, err := s.repository.First(ctx, entity, langID)
	if err != nil {
		return nil, err
	}

	if err = s.EntityInit(ctx, entity, langID); err != nil {
		return nil, err
	}
	return entity, nil
}

// Query returns the items with the specified selection condition.
func (s *service) Query(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]Entity, error) {
	items, err := s.repository.Query(ctx, cond, langID)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of items by query: %v", cond)
	}

	for i := range items {

		if err = s.EntityInit(ctx, &items[i], langID); err != nil {
			return nil, err
		}

	}
	return items, nil
}

func (s *service) Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error) {
	count, err := s.repository.Count(ctx, cond)
	if err != nil {
		return 0, errors.Wrapf(err, "Can not count a list of items by query: %v", cond)
	}
	return count, nil
}

func (s *service) Create(ctx context.Context, entity *Entity, langID uint) error {
	err := s.repository.Create(ctx, entity, langID)
	if err != nil {
		return errors.Wrapf(err, "Can not create an entity: %v", entity)
	}

	if err = s.EntityInit(ctx, entity, langID); err != nil {
		return err
	}

	return nil
}

func (s *service) Update(ctx context.Context, entity *Entity, langID uint) error {
	err := s.repository.Update(ctx, entity, langID)
	if err != nil {
		return errors.Wrapf(err, "Can not update an entity: %v", entity)
	}

	if err = s.EntityInit(ctx, entity, langID); err != nil {
		return err
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "Can not delete an entity by ID: %v", id)
	}
	return nil
}

func (s *service) EntityInit(ctx context.Context, entity *Entity, langID uint) error {

	if err := s.propertiesInit(ctx, entity, langID); err != nil {
		return err
	}

	if err := entity.propertiesValuesInit(langID); err != nil {
		return err
	}

	return nil
}

func (s *service) propertiesInit(ctx context.Context, entity *Entity, langID uint) error {
	props, rels, err := s.reference.Relation.Service.TPropertyAndRelationQuery(ctx, &selection_condition.SelectionCondition{
		Where: selection_condition.WhereCondition{
			Field:     "ID",
			Condition: selection_condition.ConditionIn,
			Value:     entity.GetPropertiesIDs(),
		},
	}, langID)
	if err != nil {
		return errors.Wrapf(err, "Can not find properties and relations for an entity: %v", entity)
	}

	entity.PropertiesValues = make(map[uint]PropertyValue, len(props))
	for i := range props {
		entity.PropertiesValues[props[i].ID] = PropertyValue{Property: &props[i]}
	}

	entity.RelationsValues = make(map[uint]RelationValue, len(rels))
	for i := range rels {
		entity.RelationsValues[rels[i].ID] = RelationValue{Relation: &rels[i]}
	}

	return nil
}

func (s *service) BindRelatedEntity(relation *entity_type.Relation, entity1 *Entity, entity2 *Entity) error {
	_, err := s.CheckBindRelatedEntity(relation, entity1, entity2)
	if err != nil {
		return err
	}

	entity1Copy := *entity1
	entity2Copy := *entity2

	err = entity1.BindRelatedEntityID(relation, entity2.ID)
	if err == nil {
		err = entity2.BindRelatedEntityID(relation, entity1.ID)
	}

	if err != nil {
		*entity1 = entity1Copy
		*entity2 = entity2Copy
	}

	return err
}

func (s *service) BindRelatedEntities(relation *entity_type.Relation, entity *Entity, relatedEntities []*Entity) error {
	_, err := s.CheckBindRelatedEntities(relation, entity, relatedEntities)
	if err != nil {
		return err
	}

	for _, relatedEntity := range relatedEntities {
		if err = s.BindRelatedEntity(relation, entity, relatedEntity); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) UnbindRelatedEntity(relation *entity_type.Relation, entity1 *Entity, entity2 *Entity) error {
	_, err := s.CheckBindRelatedEntity(relation, entity1, entity2)
	if err != nil {
		return err
	}

	entity1Copy := *entity1
	entity2Copy := *entity2

	err = entity1.UnbindRelatedEntityID(relation, entity2.ID)
	if err == nil {
		err = entity2.UnbindRelatedEntityID(relation, entity1.ID)
	}

	if err != nil {
		*entity1 = entity1Copy
		*entity2 = entity2Copy
	}

	return err
}

func (s *service) UnbindRelatedEntities(relation *entity_type.Relation, entity *Entity, relatedEntities []*Entity) error {
	_, err := s.CheckBindRelatedEntities(relation, entity, relatedEntities)
	if err != nil {
		return err
	}

	for _, relatedEntity := range relatedEntities {
		if err = s.UnbindRelatedEntity(relation, entity, relatedEntity); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) CheckBindRelatedEntity(relation *entity_type.Relation, entity1 *Entity, entity2 *Entity) (firstParamEntityIsDepended bool, err error) {
	if relation == nil {
		return false, errors.Wrapf(yaruserror.ErrNotSet, "relation is nil")
	}

	if entity1 == nil {
		return false, errors.Wrapf(yaruserror.ErrNotSet, "entity1 is nil")
	}

	if entity2 == nil {
		return false, errors.Wrapf(yaruserror.ErrNotSet, "entity2 is nil")
	}

	if entity1.EntityTypeID == relation.DependedEntityType.ID && entity2.EntityTypeID == relation.UndependedEntityType.ID {
		return true, nil
	} else if entity2.EntityTypeID == relation.DependedEntityType.ID && entity1.EntityTypeID == relation.UndependedEntityType.ID {
		return false, nil
	} else {
		return false, errors.New("wrong types")
	}

	return false, nil
}

func (s *service) CheckBindRelatedEntities(relation *entity_type.Relation, entity *Entity, relatedEntities []*Entity) (firstParamEntityIsDepended bool, err error) {
	if relation == nil {
		return false, errors.Wrapf(yaruserror.ErrNotSet, "relation is nil")
	}

	if entity == nil {
		return false, errors.Wrapf(yaruserror.ErrNotSet, "entity is nil")
	}

	if relatedEntities == nil {
		return false, errors.Wrapf(yaruserror.ErrNotSet, "relatedEntities is nil")
	}

	if len(relatedEntities) == 0 {
		return false, errors.Wrapf(yaruserror.ErrNotSet, "relatedEntities is empty")
	}

	if entity.EntityTypeID == relation.DependedEntityType.ID {
		for _, e := range relatedEntities {
			if e.EntityTypeID != relation.UndependedEntityType.ID {
				return false, errors.New("wrong types")
			}
		}
	} else if entity.EntityTypeID == relation.UndependedEntityType.ID {
		for _, e := range relatedEntities {
			if e.EntityTypeID != relation.DependedEntityType.ID {
				return false, errors.New("wrong types")
			}
		}
	} else {
		return false, errors.New("wrong types")
	}

	return false, nil
}
