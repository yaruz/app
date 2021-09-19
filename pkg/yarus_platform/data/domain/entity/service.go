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
	Query(ctx context.Context, query *selection_condition.SelectionCondition) ([]Entity, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *Entity) error
	Update(ctx context.Context, entity *Entity) error
	Save(ctx context.Context, entity *Entity) error
	Delete(ctx context.Context, id uint) error
	EntityInit(ctx context.Context, entity *Entity, langID uint) error
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

// Query returns the items with the specified selection condition.
func (s *service) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]Entity, error) {
	items, err := s.repository.Query(ctx, cond)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of items by query: %v", cond)
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

func (s *service) Create(ctx context.Context, entity *Entity) error {
	err := s.repository.Create(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not create an entity: %v", entity)
	}
	return nil
}

func (s *service) Update(ctx context.Context, entity *Entity) error {
	err := s.repository.Update(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not update an entity: %v", entity)
	}
	return nil
}

func (s *service) Save(ctx context.Context, entity *Entity) error {
	err := s.repository.Save(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not save an entity: %v", entity)
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
	return nil
}

func (s *service) propertiesInit(ctx context.Context, entity *Entity) error {
	return s.propsInit(ctx, entity, func(ctx context.Context, ids []interface{}) ([]property.Property, []entity_type.Relation, error) {
		return s.reference.Relation.Service.PropertyAndRelationQuery(ctx, &selection_condition.SelectionCondition{
			Where: selection_condition.WhereCondition{
				Field:     "ID",
				Condition: "in",
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
				Condition: "in",
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
	for _, prop := range props {
		entity.PropertiesValues[prop.ID] = PropertyValue{Property: prop}
	}

	entity.RelationsValues = make(map[uint]RelationValue, len(rels))
	for _, rel := range rels {
		entity.RelationsValues[rel.ID] = RelationValue{Relation: rel}
	}

	return nil
}

func (s *service) tPropertiesValuesInit(ctx context.Context, entity *Entity, langID uint) error {

	for id, val := range entity.PropertiesValuesMap {

		_, propOk := entity.PropertiesValues[id]
		rel, relOk := entity.RelationsValues[id]

		switch {
		case relOk:
			entitiesIDs, ok := val.([]uint)
			if !ok {
				return errors.Errorf("Can not cast value of relation into a []uint. ID = %v; Val = %v.", id, val)
			}
			rel.Value = entitiesIDs
			entity.RelationsValues[id] = rel
		case propOk:
		default:
			return errors.Errorf("Property ID = %v not found.", id)
		}

	}

	return nil
}
