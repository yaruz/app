package sn_account

import (
	"context"

	"github.com/minipkg/selection_condition"
	"github.com/pkg/errors"

	"github.com/minipkg/log"
)

// IService encapsulates usecase logic for SnAccount.
type IService interface {
	New(ctx context.Context) (*SnAccount, error)
	Get(ctx context.Context, id uint, langID uint) (*SnAccount, error)
	Query(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) ([]SnAccount, error)
	First(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (*SnAccount, error)
	Count(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (uint, error)
	Create(ctx context.Context, obj *SnAccount, langID uint) error
	Update(ctx context.Context, obj *SnAccount, langID uint) error
	Delete(ctx context.Context, id uint) error
}

type service struct {
	//Domain     Domain
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
func (s service) defaultConditions() *selection_condition.SelectionCondition {
	return &selection_condition.SelectionCondition{}
}

func (s *service) New(ctx context.Context) (*SnAccount, error) {
	return s.repository.New(ctx)
}

// Get returns the entity with the specified ID.
func (s *service) Get(ctx context.Context, id uint, langID uint) (*SnAccount, error) {
	entity, err := s.repository.Get(ctx, id, langID)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not get a %v by id: %v", EntityType, id)
	}
	return entity, nil
}

// Query returns the items with the specified offset and limit.
func (s *service) Query(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) ([]SnAccount, error) {
	return s.repository.Query(ctx, condition, langID)
}

func (s *service) First(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (*SnAccount, error) {
	return s.repository.First(ctx, condition, langID)
}

// Count returns the number of items.
func (s *service) Count(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (uint, error) {
	return s.repository.Count(ctx, condition, langID)
}

func (s *service) Create(ctx context.Context, obj *SnAccount, langID uint) error {
	return s.repository.Create(ctx, obj, langID)
}

func (s *service) Update(ctx context.Context, obj *SnAccount, langID uint) error {
	return s.repository.Update(ctx, obj, langID)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repository.Delete(ctx, id)
}
