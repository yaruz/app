package user

import (
	"context"
	"github.com/minipkg/selection_condition"

	"github.com/pkg/errors"

	"github.com/minipkg/log"
)

// IService encapsulates usecase logic for user.
type IService interface {
	NewEntity() *User
	Get(ctx context.Context, id uint) (*User, error)
	//Query(ctx context.Context, offset, limit uint) ([]User, error)
	//List(ctx context.Context) ([]User, error)
	//Count(ctx context.Context) (uint, error)
	Create(ctx context.Context, entity *User) error
	//Update(ctx context.Context, id string, input *User) (*User, error)
	//Delete(ctx context.Context, id string) (error)
	First(ctx context.Context, user *User) (*User, error)
}

type service struct {
	//Domain     Domain
	logger log.ILogger
	repo   Repository
}

// NewService creates a new service.
func NewService(logger log.ILogger, repo Repository) IService {
	s := &service{
		logger: logger,
		repo:   repo,
	}
	repo.SetDefaultConditions(s.defaultConditions())
	return s
}

// Defaults returns defaults params
func (s service) defaultConditions() *selection_condition.SelectionCondition {
	return &selection_condition.SelectionCondition{}
}

func (s service) NewEntity() *User {
	return &User{}
}

// Get returns the entity with the specified ID.
func (s service) Get(ctx context.Context, id uint) (*User, error) {
	entity, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not get a user by id: %v", id)
	}
	return entity, nil
}

/*
// Count returns the number of items.
func (s service) Count(ctx context.Context) (uint, error) {
	return s.repo.Count(ctx)
}*/

// Query returns the items with the specified offset and limit.
/*func (s service) Query(ctx context.Context, offset, limit uint) ([]User, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of users by ctx")
	}
	return items, nil
}

// List returns the items list.
func (s service) List(ctx context.Context) ([]User, error) {
	items, err := s.repo.Query(ctx, 0, MaxLIstLimit)
	if err != nil {
		return nil, errors.Wrapf(err, "Can not find a list of users by ctx")
	}
	return items, nil
}*/

func (s service) Create(ctx context.Context, entity *User) error {
	return s.repo.Create(ctx, entity)
}

func (s service) First(ctx context.Context, user *User) (*User, error) {
	return s.repo.First(ctx, user)
}
