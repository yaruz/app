package entity_type

import (
	"context"

	"github.com/minipkg/selection_condition"

	"github.com/pkg/errors"

	"github.com/minipkg/log"
)

// IService encapsulates usecase logic.
type IService interface {
	NewEntity() *EntityType
	Get(ctx context.Context, id uint) (*EntityType, error)
	First(ctx context.Context, entity *EntityType) (*EntityType, error)
	Query(ctx context.Context, query *selection_condition.SelectionCondition) ([]EntityType, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *EntityType) error
	Update(ctx context.Context, entity *EntityType) error
	Delete(ctx context.Context, entity *EntityType) error
	TGet(ctx context.Context, id uint, langID uint) (*EntityType, error)
	TFirst(ctx context.Context, entity *EntityType, langID uint) (*EntityType, error)
	TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]EntityType, error)
	TCreate(ctx context.Context, entity *EntityType, langID uint) (err error)
	TUpdate(ctx context.Context, entity *EntityType, langID uint) (err error)
	BindProperty(ctx context.Context, id uint, propertyID uint) error
	UnbindProperty(ctx context.Context, id uint, propertyID uint) error
	UnbindAllProperty(ctx context.Context, id uint) error
}

type service struct {
	logger          log.ILogger
	repository      Repository
	relationService RelationService
}

var _ IService = (*service)(nil)

// NewService creates a new service.
func NewService(logger log.ILogger, repo Repository, relationService RelationService) IService {
	s := &service{
		logger:          logger,
		repository:      repo,
		relationService: relationService,
	}
	repo.SetDefaultConditions(s.defaultConditions())
	return s
}

// Defaults returns defaults params
func (s *service) defaultConditions() *selection_condition.SelectionCondition {
	return &selection_condition.SelectionCondition{}
}

func (s *service) NewEntity() *EntityType {
	return New()
}

func (s *service) initPropertiesAndRelations(ctx context.Context, entity *EntityType) (err error) {
	(*entity).Properties, (*entity).Relations, err = s.relationService.GetPropertiesAndRelationsByEntityTypeID(ctx, (*entity).ID)
	return err
}

func (s *service) tInitPropertiesAndRelations(ctx context.Context, entity *EntityType, langID uint) (err error) {
	(*entity).Properties, (*entity).Relations, err = s.relationService.TGetPropertiesAndRelationsByEntityTypeID(ctx, (*entity).ID, langID)
	return err
}

// Get returns the entity with the specified ID.
func (s *service) Get(ctx context.Context, id uint) (*EntityType, error) {
	entity, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return entity, s.initPropertiesAndRelations(ctx, entity)
}

func (s *service) First(ctx context.Context, entity *EntityType) (*EntityType, error) {
	entity, err := s.repository.First(ctx, entity)
	if err != nil {
		return nil, err
	}
	return entity, s.initPropertiesAndRelations(ctx, entity)
}

func (s *service) TGet(ctx context.Context, id uint, langID uint) (*EntityType, error) {
	entity, err := s.repository.TGet(ctx, id, langID)
	if err != nil {
		return nil, err
	}
	return entity, s.tInitPropertiesAndRelations(ctx, entity, langID)
}

func (s *service) TFirst(ctx context.Context, entity *EntityType, langID uint) (*EntityType, error) {
	entity, err := s.repository.TFirst(ctx, entity, langID)
	if err != nil {
		return nil, err
	}
	return entity, s.tInitPropertiesAndRelations(ctx, entity, langID)
}

// Query returns the items with the specified selection condition.
func (s *service) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]EntityType, error) {
	items, err := s.repository.Query(ctx, cond)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of items by query: %v", cond)
	}

	for i := range items {
		if err = s.initPropertiesAndRelations(ctx, &items[i]); err != nil {
			return nil, err
		}
	}
	return items, nil
}

func (s *service) TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]EntityType, error) {
	items, err := s.repository.TQuery(ctx, cond, langID)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of items by query: %v", cond)
	}

	for i := range items {
		if err = s.tInitPropertiesAndRelations(ctx, &items[i], langID); err != nil {
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

func (s *service) Create(ctx context.Context, entity *EntityType) error {
	err := entity.Validate()
	if err != nil {
		return errors.Wrapf(err, "Validation error: %v", err)
	}

	err = s.repository.Create(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not create an entity: %v", entity)
	}
	return nil
}

func (s *service) TCreate(ctx context.Context, entity *EntityType, langID uint) error {
	err := entity.Validate()
	if err != nil {
		return errors.Wrapf(err, "Validation error: %v", err)
	}

	err = s.repository.TCreate(ctx, entity, langID)
	if err != nil {
		return errors.Wrapf(err, "Can not create an entity: %v", entity)
	}
	return nil
}

func (s *service) Update(ctx context.Context, entity *EntityType) error {
	err := entity.Validate()
	if err != nil {
		return errors.Wrapf(err, "Validation error: %v", err)
	}

	err = s.repository.Update(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not update an entity: %v", entity)
	}
	return nil
}

func (s *service) TUpdate(ctx context.Context, entity *EntityType, langID uint) error {
	err := entity.Validate()
	if err != nil {
		return errors.Wrapf(err, "Validation error: %v", err)
	}

	err = s.repository.TUpdate(ctx, entity, langID)
	if err != nil {
		return errors.Wrapf(err, "Can not update an entity: %v", entity)
	}
	return nil
}

func (s *service) Delete(ctx context.Context, entity *EntityType) error {
	err := s.repository.Delete(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not delete an entity: %v", entity)
	}
	return nil
}

func (s *service) BindProperty(ctx context.Context, id uint, propertyID uint) error {
	err := s.repository.BindProperty(ctx, id, propertyID)
	if err != nil {
		return errors.Wrapf(err, "Can not for an entity ID = %v bind an property ID = %v", id, propertyID)
	}
	return nil
}

func (s *service) UnbindProperty(ctx context.Context, id uint, propertyID uint) error {
	err := s.repository.UnbindProperty(ctx, id, propertyID)
	if err != nil {
		return errors.Wrapf(err, "Can not for an entity ID = %v unbind an property ID = %v", id, propertyID)
	}
	return nil
}

func (s *service) UnbindAllProperty(ctx context.Context, id uint) error {
	err := s.repository.UnbindAllProperty(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "Can not for an entity ID = %v unbind all properties", id)
	}
	return nil
}
