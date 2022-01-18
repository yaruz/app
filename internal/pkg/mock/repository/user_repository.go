package repository

import (
	"context"

	"github.com/minipkg/selection_condition"
	"github.com/stretchr/testify/mock"

	"github.com/yaruz/app/internal/domain/user"
)

// UserRepository is a mock for UserRepository
type UserRepository struct {
	mock.Mock
}

var _ user.Repository = (*UserRepository)(nil)

func (m UserRepository) SetDefaultConditions(conditions *selection_condition.SelectionCondition) {}

func (m UserRepository) New(ctx context.Context) (*user.User, error) {
	return &user.User{}, nil
}

func (m UserRepository) Count(ctx context.Context, condition *selection_condition.SelectionCondition, langID uint) (uint, error) {
	return 0, nil
}

func (m UserRepository) Update(ctx context.Context, obj *user.User, langID uint) error {
	return nil
}

func (m UserRepository) Delete(ctx context.Context, id uint) error {
	return nil
}

func (m UserRepository) Get(a0 context.Context, a1 uint, langID uint) (*user.User, error) {
	ret := m.Called(a0, a1)

	var r0 *user.User
	if rf, ok := ret.Get(0).(func(context.Context, uint) *user.User); ok {
		r0 = rf(a0, a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(a0, a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m UserRepository) First(a0 context.Context, a1 *selection_condition.SelectionCondition, a2 uint) (*user.User, error) {
	ret := m.Called(a0, a1)

	var r0 *user.User
	if rf, ok := ret.Get(0).(func(context.Context, *selection_condition.SelectionCondition, uint) *user.User); ok {
		r0 = rf(a0, a1, a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *selection_condition.SelectionCondition, uint) error); ok {
		r1 = rf(a0, a1, a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m UserRepository) Query(a0 context.Context, a1 *selection_condition.SelectionCondition, a2 uint) ([]user.User, error) {
	ret := m.Called(a0, a1, a2)

	var r0 []user.User
	if rf, ok := ret.Get(0).(func(context.Context, *selection_condition.SelectionCondition, uint) []user.User); ok {
		r0 = rf(a0, a1, a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]user.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *selection_condition.SelectionCondition, uint) error); ok {
		r1 = rf(a0, a1, a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m UserRepository) Create(a0 context.Context, a1 *user.User, langID uint) error {
	ret := m.Called(a0, a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *user.User) error); ok {
		r0 = rf(a0, a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
