package property_group

import (
	"context"

	"github.com/minipkg/selection_condition"

	"github.com/pkg/errors"

	"github.com/minipkg/log"
)

// IService encapsulates usecase logic.
type IService interface {
	NewEntity() *PropertyGroup
	Get(ctx context.Context, id uint) (*PropertyGroup, error)
	First(ctx context.Context, entity *PropertyGroup) (*PropertyGroup, error)
	Query(ctx context.Context, query *selection_condition.SelectionCondition) ([]PropertyGroup, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (int64, error)
	Create(ctx context.Context, entity *PropertyGroup) error
	Update(ctx context.Context, entity *PropertyGroup) error
	Delete(ctx context.Context, entity *PropertyGroup) error
	TGet(ctx context.Context, id uint, langID uint) (*PropertyGroup, error)
	TFirst(ctx context.Context, entity *PropertyGroup, langID uint) (*PropertyGroup, error)
	TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]PropertyGroup, error)
	TCreate(ctx context.Context, entity *PropertyGroup, langID uint) (err error)
	TUpdate(ctx context.Context, entity *PropertyGroup, langID uint) (err error)
}

type service struct {
	logger     log.ILogger
	repository Repository
}

var _ IService = (*service)(nil)

// NewService creates a new service.
func NewService(logger log.ILogger, repo Repository) IService {
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

func (s *service) NewEntity() *PropertyGroup {
	return New()
}

// Get returns the entity with the specified ID.
func (s *service) Get(ctx context.Context, id uint) (*PropertyGroup, error) {
	entity, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *service) First(ctx context.Context, entity *PropertyGroup) (*PropertyGroup, error) {
	entity, err := s.repository.First(ctx, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *service) TGet(ctx context.Context, id uint, langID uint) (*PropertyGroup, error) {
	entity, err := s.repository.TGet(ctx, id, langID)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *service) TFirst(ctx context.Context, entity *PropertyGroup, langID uint) (*PropertyGroup, error) {
	entity, err := s.repository.TFirst(ctx, entity, langID)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Query returns the items with the specified selection condition.
func (s *service) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]PropertyGroup, error) {
	items, err := s.repository.Query(ctx, cond)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of items by query: %v", cond)
	}
	return items, nil
}

func (s *service) TQuery(ctx context.Context, cond *selection_condition.SelectionCondition, langID uint) ([]PropertyGroup, error) {
	items, err := s.repository.TQuery(ctx, cond, langID)
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

func (s *service) Create(ctx context.Context, entity *PropertyGroup) error {
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

func (s *service) TCreate(ctx context.Context, entity *PropertyGroup, langID uint) error {
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

func (s *service) Update(ctx context.Context, entity *PropertyGroup) error {
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

func (s *service) TUpdate(ctx context.Context, entity *PropertyGroup, langID uint) error {
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

func (s *service) Delete(ctx context.Context, entity *PropertyGroup) error {
	err := s.repository.Delete(ctx, entity)
	if err != nil {
		return errors.Wrapf(err, "Can not delete an entity: %v", entity)
	}
	return nil
}
