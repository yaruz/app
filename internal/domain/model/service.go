package model

import (
	"context"
	"github.com/minipkg/selection_condition"

	"github.com/pkg/errors"

	"github.com/minipkg/log"
)

// IService encapsulates usecase logic for user.
type IService interface {
	NewEntity() *Model
	Get(ctx context.Context, id uint) (*Model, error)
	Query(ctx context.Context, query *selection_condition.SelectionCondition) ([]Model, error)
	Count(ctx context.Context, cond *selection_condition.SelectionCondition) (uint, error)
}

type service struct {
	logger     log.ILogger
	repository Repository
}

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
func (s service) defaultConditions() *selection_condition.SelectionCondition {
	return &selection_condition.SelectionCondition{}
}

func (s service) NewEntity() *Model {
	return New()
}

// Get returns the entity with the specified ID.
func (s service) Get(ctx context.Context, id uint) (*Model, error) {
	entity, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// Query returns the items with the specified offset and limit.
func (s service) Query(ctx context.Context, query *selection_condition.SelectionCondition) ([]Model, error) {
	items, err := s.repository.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of Model by query: %v", query)
	}
	return items, nil
}

func (s service) Count(ctx context.Context, query *selection_condition.SelectionCondition) (uint, error) {
	count, err := s.repository.Count(ctx, query)
	if err != nil {
		return 0, errors.Wrapf(err, "Can not find a list of Model by query: %v", query)
	}
	return count, nil
}
