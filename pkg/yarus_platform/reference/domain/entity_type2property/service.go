package entity_type2property

import (
	"context"

	"github.com/minipkg/selection_condition"

	"github.com/pkg/errors"

	"github.com/minipkg/log"
)

// IService encapsulates usecase logic.
type IService interface {
	NewEntity() *EntityType2Property
	Get(ctx context.Context, id uint) (*EntityType2Property, error)
	Query(ctx context.Context, query *selection_condition.SelectionCondition) ([]EntityType2Property, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *EntityType2Property) error
	Delete(ctx context.Context, entity *EntityType2Property) error
}

type service struct {
	logger     log.Logger
	repository Repository
}

var _ IService = (*service)(nil)

// NewService creates a new service.
func NewService(logger log.Logger, repo Repository) IService {
	s := &service{
		logger:     logger,
		repository: repo,
	}
	repo.SetDefaultConditions(s.defaultConditions())
	return s
}

// Defaults returns defaults params
func (s *service) defaultConditions() *selection_condition.SelectionCondition {
	return &selection_condition.SelectionCondition{}
}

func (s *service) NewEntity() *EntityType2Property {
	return New()
}

// Get returns the entity with the specified ID.
func (s *service) Get(ctx context.Context, id uint) (*EntityType2Property, error) {
	entity, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Query returns the items with the specified selection condition.
func (s *service) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]EntityType2Property, error) {
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

func (s *service) Create(ctx context.Context, entity *EntityType2Property) error {
	err := s.repository.Create(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not create an entity: %v", entity)
	}
	return nil
}

func (s *service) Delete(ctx context.Context, entity *EntityType2Property) error {
	err := s.repository.Delete(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not delete an entity: %#v", entity)
	}
	return nil
}
