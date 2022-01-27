package redis

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/minipkg/selection_condition"
	"github.com/yaruz/app/internal/pkg/session"

	mock2 "github.com/minipkg/db/redis/mock"

	"github.com/elliotchance/redismock/v8"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	dbredis "github.com/minipkg/db/redis"

	"github.com/yaruz/app/internal/domain/user"
)

const sessionLifeTimeInHours = 1

type SessionRepositoryTestSuite struct {
	//	for all tests
	suite.Suite
	user    *user.User
	session *session.Session
	//	only for each individual test
	ctx        context.Context
	mock       *redismock.ClientMock
	repository *SessionRepository
}

type userRepoMock struct {
	mock.Mock
	user *user.User
}

var _ user.Repository = (*userRepoMock)(nil)

func (m *userRepoMock) SetDefaultConditions(conditions *selection_condition.SelectionCondition) {

}

func (m *userRepoMock) Get(ctx context.Context, id uint) (*user.User, error) {
	return m.user, nil
}

func (m *userRepoMock) Query(ctx context.Context, cond *selection_condition.SelectionCondition) ([]user.User, error) {
	return []user.User{*m.user}, nil
}

func (m *userRepoMock) Create(ctx context.Context, entity *user.User) error {
	return nil
}

func (m *userRepoMock) First(ctx context.Context, user *user.User) (*user.User, error) {
	return m.user, nil
}

func (s *SessionRepositoryTestSuite) SetupSuite() {
	var passhash, _ = hex.DecodeString("3a73acfdb534ddded4c0109383ee3e5a66314113d1ff691aaf4b3ee073c8fc2edd06d48f0555ec3783f4c479994e3eee3433734c29b05f08be0e9739b956b88d8fe872bd0a0942214e94fd4001e757fa3b66a2b9925de2e800c55ef49baa4c03")

	s.user = &user.User{
		ID:        1,
		Name:      "demo1",
		Passhash:  string(passhash),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.session = &session.Session{
		UserID: s.user.ID,
		User:   *s.user,
	}
}

func (s *SessionRepositoryTestSuite) SetupTest() {
	var db *dbredis.DB
	var err error
	require := require.New(s.T())
	s.session.Data = session.Data{}

	db, s.mock, err = mock2.New()
	require.NoError(err)

	s.repository, err = NewSessionRepository(db, sessionLifeTimeInHours, &userRepoMock{user: s.user})
	require.NoError(err)
}

func TestSessionRepository(t *testing.T) {
	suite.Run(t, new(SessionRepositoryTestSuite))
}

func (s *SessionRepositoryTestSuite) TestNewEntity() {
	assert := assert.New(s.T())
	require := require.New(s.T())

	res, err := s.repository.NewEntity(s.ctx, s.user.ID)
	require.NoError(err)

	assert.Equalf(*s.session, *res, "The two objects should be the same. Expected: %v; have got: %v", *s.session, *res)
}

func (s *SessionRepositoryTestSuite) TestGet() {
	assert := assert.New(s.T())
	require := require.New(s.T())

	jsonSess, err := s.session.MarshalBinary()
	require.NoError(err)

	s.mock.On("Get", s.ctx, s.repository.Key(s.user.ID)).
		Return(redis.NewStringResult(string(jsonSess), nil))

	res, err := s.repository.Get(s.ctx, s.user.ID)
	require.NoError(err)

	jsonRes, err := res.MarshalBinary()
	require.NoError(err)

	assert.Equalf(jsonSess, jsonRes, "The two objects should be the same. Expected: %v; have got: %v", jsonSess, jsonRes)
}

func (s *SessionRepositoryTestSuite) TestCreate() {
	var err error
	require := require.New(s.T())

	s.mock.On("Set", s.ctx, s.repository.Key(s.user.ID), s.session, s.repository.SessionLifeTime).
		Return(redis.NewStatusResult("", nil))

	err = s.repository.Create(s.ctx, s.session)
	require.NoError(err)
}

func (s *SessionRepositoryTestSuite) TestUpdate() {
	var err error
	require := require.New(s.T())

	s.mock.On("Set", s.ctx, s.repository.Key(s.user.ID), s.session, s.repository.SessionLifeTime).
		Return(redis.NewStatusResult("", nil))

	err = s.repository.Update(s.ctx, s.session)
	require.NoError(err)
}

func (s *SessionRepositoryTestSuite) TestDelete() {
	var err error
	require := require.New(s.T())

	s.mock.On("Del", s.ctx, []string{s.repository.Key(s.user.ID)}).
		Return(redis.NewIntResult(1, nil))

	err = s.repository.Delete(s.ctx, s.session)
	require.NoError(err)
}

func (s *SessionRepositoryTestSuite) TestSetData() {
	var err error
	require := require.New(s.T())
	s.session.Data.UserID = 1
	s.session.Data.UserName = "User1"
	s.session.Data.ExpirationTokenTime = time.Now()

	s.mock.On("Set", s.ctx, s.repository.Key(s.user.ID), s.session, s.repository.SessionLifeTime).
		Return(redis.NewStatusResult("", nil))

	err = s.repository.SetData(s.session, s.session.Data)
	require.NoError(err)
}

func (s *SessionRepositoryTestSuite) TestGetData() {
	require := require.New(s.T())
	s.session.Data.UserID = 1
	s.session.Data.UserName = "User1"
	s.session.Data.ExpirationTokenTime = time.Now()

	s.mock.On("Set", s.ctx, s.repository.Key(s.user.ID), *s.session, s.repository.SessionLifeTime).
		Return(redis.NewStatusResult("", nil))

	res := s.repository.GetData(s.session)
	require.Equalf(s.session.Data, res, "The two objects should be the same. Expected: %v; have got: %v", s.session.Data, res)
}
