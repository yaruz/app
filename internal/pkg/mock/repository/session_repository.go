package repository

import (
	"context"
	"github.com/minipkg/selection_condition"

	"github.com/stretchr/testify/mock"

	"github.com/yaruz/app/internal/pkg/auth"
	"github.com/yaruz/app/internal/pkg/session"
)

// SessionRepository is a mock for SessionRepository
type SessionRepository struct {
	mock.Mock
}

var _ auth.SessionRepository = (*SessionRepository)(nil)

func (m SessionRepository) SetDefaultConditions(defaultConditions selection_condition.SelectionCondition) {
}

func (m SessionRepository) NewEntity(a0 context.Context, a1 uint) (*session.Session, error) {
	ret := m.Called(a0, a1)

	var r0 *session.Session
	if rf, ok := ret.Get(0).(func(context.Context, uint) *session.Session); ok {
		r0 = rf(a0, a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*session.Session)
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

func (m SessionRepository) Get(a0 context.Context, a1 uint) (*session.Session, error) {
	ret := m.Called(a0, a1)

	var r0 *session.Session
	if rf, ok := ret.Get(0).(func(context.Context, uint) *session.Session); ok {
		r0 = rf(a0, a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*session.Session)
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

func (m SessionRepository) Create(a0 context.Context, a1 *session.Session) error {
	ret := m.Called(a0, a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *session.Session) error); ok {
		r0 = rf(a0, a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m SessionRepository) Update(a0 context.Context, a1 *session.Session) error {
	ret := m.Called(a0, a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *session.Session) error); ok {
		r0 = rf(a0, a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m SessionRepository) Save(a0 *session.Session) error {
	ret := m.Called(a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*session.Session) error); ok {
		r0 = rf(a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m SessionRepository) Delete(a0 context.Context, a1 *session.Session) error {
	ret := m.Called(a0, a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *session.Session) error); ok {
		r0 = rf(a0, a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (m SessionRepository) GetData(a0 *session.Session) session.Data {
	ret := m.Called(a0)

	var r0 session.Data
	if rf, ok := ret.Get(0).(func(*session.Session) session.Data); ok {
		r0 = rf(a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(session.Data)
		}
	}

	return r0
}

func (m SessionRepository) SetData(a0 *session.Session, a1 session.Data) error {
	ret := m.Called(a0, a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*session.Session, session.Data) error); ok {
		r0 = rf(a0, a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
