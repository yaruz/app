package entity

import (
	"context"

	"github.com/pkg/errors"

	"github.com/minipkg/log"
	"github.com/minipkg/selection_condition"

	"github.com/yaruz/app/pkg/yarus_platform/reference"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/entity_type"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/property"
	"github.com/yaruz/app/pkg/yarus_platform/search"
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
	EntitySetPropertyValue(ctx context.Context, entity *Entity, property *property.Property, value interface{}, langID uint) error
	EntitySetRelationValue(ctx context.Context, entity *Entity, relation *entity_type.Relation, value []uint) error
	EntityDeletePropertyValue(ctx context.Context, entity *Entity, propertyID uint)
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
	repo.SetDefaultConditions(s.defaultConditions())
	return s
}

// Defaults returns defaults params
func (s *service) defaultConditions() *selection_condition.SelectionCondition {
	return &selection_condition.SelectionCondition{}
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

	if err := s.tPropertiesInit(ctx, entity, langID); err != nil {
		return err
	}

	if err := s.tPropertiesValuesInit(ctx, entity, langID); err != nil {
		return err
	}

	return nil
}

func (s *service) propertiesInit(ctx context.Context, entity *Entity) error {
	return s.propsInit(ctx, entity, func(ctx context.Context, ids []interface{}) ([]property.Property, []entity_type.Relation, error) {
		return s.reference.Relation.Service.PropertyAndRelationQuery(ctx, &selection_condition.SelectionCondition{
			Where: selection_condition.WhereCondition{
				Field:     "ID",
				Condition: selection_condition.ConditionIn,
				Value:     ids,
			},
		})
	})
}

func (s *service) tPropertiesInit(ctx context.Context, entity *Entity, langID uint) error {
	return s.propsInit(ctx, entity, func(ctx context.Context, ids []interface{}) ([]property.Property, []entity_type.Relation, error) {
		return s.reference.Relation.Service.TPropertyAndRelationQuery(ctx, &selection_condition.SelectionCondition{
			Where: selection_condition.WhereCondition{
				Field:     "ID",
				Condition: selection_condition.ConditionIn,
				Value:     ids,
			},
		}, langID)
	})
}

func (s *service) propsInit(ctx context.Context, entity *Entity, propertyAndRelationQuery func(ctx context.Context, ids []interface{}) ([]property.Property, []entity_type.Relation, error)) error {
	ids := make([]interface{}, 0, len(entity.PropertiesValuesMap))
	for id := range entity.PropertiesValuesMap {
		ids = append(ids, id)
	}

	props, rels, err := propertyAndRelationQuery(ctx, ids)
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

// Инициализация PropertiesValues и RelationsValues из PropertiesValuesMap
// Запускаем после запуска tPropertiesInit
func (s *service) tPropertiesValuesInit(ctx context.Context, entity *Entity, langID uint) error {

	for id, val := range entity.PropertiesValuesMap {

		prop, propOk := entity.PropertiesValues[id]
		rel, relOk := entity.RelationsValues[id]

		switch {
		case relOk:
			if err := rel.SetValueByInterface(val); err != nil {
				return errors.Wrapf(err, "Can not set value to PropertyValue. Property ID = %v; Value = %v.", id, val)
			}
			entity.RelationsValues[id] = rel
		case propOk:
			if err := prop.SetValue(val, langID); err != nil {
				return errors.Wrapf(err, "Can not set value to PropertyValue. Property ID = %v; Value = %v.", id, val)
			}
			entity.PropertiesValues[id] = prop
		default:
			return errors.Errorf("Property ID = %v not found.", id)
		}

	}

	return nil
}

// value - значение, не ссылка
func (s *service) EntitySetPropertyValue(ctx context.Context, entity *Entity, property *property.Property, value interface{}, langID uint) error {
	propertyValue, err := newPropertyValue(property, value, langID)
	if err != nil {
		return err
	}
	entity.setPropertyValue(propertyValue)
	return nil
}

// value - значение, не ссылка, []uint - IDs связанных сущностей
func (s *service) EntitySetRelationValue(ctx context.Context, entity *Entity, relation *entity_type.Relation, value []uint) error {
	relationValue, err := newRelationValue(relation, value)
	if err != nil {
		return err
	}
	entity.setRelationValue(relationValue)
	return nil
}

// Удаляет как значения свойств, так и значения связей
func (s *service) EntityDeletePropertyValue(ctx context.Context, entity *Entity, propertyID uint) {
	entity.DeletePropertyValue(propertyID)
}

func (s *service) BindRelatedEntity(ctx context.Context, entity *Entity, relation *entity_type.Relation, entityID uint) {
}

func (s *service) BindRelatedEntities(ctx context.Context, entity *Entity, relation *entity_type.Relation, entityIDs []uint) {
}

func (s *service) UnbindRelatedEntity(ctx context.Context, entity *Entity, relation *entity_type.Relation, entityID uint) {
}

func (s *service) UnbindRelatedEntities(ctx context.Context, entity *Entity, relation *entity_type.Relation, entityIDs []uint) {
}
