package gorm

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/minipkg/db/gorm/mock"
	"github.com/minipkg/go-app-common/db/pg"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/yaruz/app/internal/pkg/config"

	"github.com/minipkg/log"

	"github.com/yaruz/app/internal/domain/user"
)

const pkgName = "pg"

type UserRepositoryTestSuite struct {
	//	for all tests
	suite.Suite
	cfg    *config.Configuration
	logger *log.Logger
	user   *user.User
	//	only for each individual test
	ctx        context.Context
	mock       sqlmock.Sqlmock
	repository user.Repository
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	var passhash, _ = hex.DecodeString("3a73acfdb534ddded4c0109383ee3e5a66314113d1ff691aaf4b3ee073c8fc2edd06d48f0555ec3783f4c479994e3eee3433734c29b05f08be0e9739b956b88d8fe872bd0a0942214e94fd4001e757fa3b66a2b9925de2e800c55ef49baa4c03")
	var err error

	s.cfg = config.Get4UnitTest("UserRepository")

	s.logger, err = log.New(s.cfg.Log)
	require.NoError(s.T(), err)

	s.user = &user.User{
		ID:        1,
		Name:      "demo1",
		Passhash:  string(passhash),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (s *UserRepositoryTestSuite) SetupTest() {
	var ok bool
	var pgDB pg.IDB
	var PgMock *sqlmock.Sqlmock
	var err error
	require := require.New(s.T())
	s.ctx = context.Background()

	pgDB, PgMock, err = mock.New(s.cfg.DB.Identity, s.logger)
	require.NoError(err)
	s.mock = *PgMock

	r, err := GetRepository(s.logger, pgDB, user.EntityName)
	require.NoError(err)

	s.repository, ok = r.(user.Repository)
	require.Truef(ok, "Can not cast DB repository for entity %q to %vRepository. Repo: %v", user.EntityName, user.EntityName, r)
}

func (s *UserRepositoryTestSuite) AfterTest(_, _ string) {
	err := s.mock.ExpectationsWereMet()
	assert.Nil(s.T(), err, "there were unfulfilled expectations: %s", err)
	//require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) TestGet() {
	assert := assert.New(s.T())

	sql := fmt.Sprintf(`SELECT \* FROM "user".*?"user"\."id" = %v.*?LIMIT 1`, s.user.ID)
	rows := sqlmock.NewRows([]string{"id", "name", "passhash", "created_at", "updated_at", "deleted_at"}).AddRow(s.user.ID, s.user.Name, s.user.Passhash, s.user.CreatedAt, s.user.UpdatedAt, s.user.DeletedAt)
	s.mock.ExpectQuery(sql).WillReturnRows(rows)

	res, err := s.repository.Get(s.ctx, s.user.ID)
	assert.Nil(err)

	assert.Equalf(*s.user, *res, "The two objects should be the same. Expected: %v; have got: %v", *s.user, *res)
}

func (s *UserRepositoryTestSuite) TestCreate() {
	assert := assert.New(s.T())

	s.mock.ExpectBegin()

	sql := fmt.Sprintf(`INSERT INTO "user".*?VALUES \(\$1,\$2,\$3,\$4,\$5\).*?RETURNING "user"\."id"`)
	rows := sqlmock.NewRows([]string{"id"}).AddRow(s.user.ID)
	s.mock.ExpectQuery(sql).WithArgs(s.user.Name, s.user.Passhash, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).WillReturnRows(rows)

	s.mock.ExpectCommit()

	user := user.New()
	user.Name = s.user.Name
	user.Passhash = s.user.Passhash

	err := s.repository.Create(s.ctx, user)
	assert.Nil(err)
}

func (s *UserRepositoryTestSuite) TestFirst() {
	assert := assert.New(s.T())

	sql := fmt.Sprintf(`SELECT \* FROM "user".*?"user"\."name" = \$1.*?LIMIT 1`)
	rows := sqlmock.NewRows([]string{"id", "name", "passhash", "created_at", "updated_at", "deleted_at"}).AddRow(s.user.ID, s.user.Name, s.user.Passhash, s.user.CreatedAt, s.user.UpdatedAt, s.user.DeletedAt)
	s.mock.ExpectQuery(sql).WillReturnRows(rows).WithArgs(s.user.Name)

	user := user.New()
	user.Name = s.user.Name

	res, err := s.repository.First(s.ctx, user)
	assert.Nil(err)

	assert.Equalf(*s.user, *res, "The two objects should be the same. Expected: %v; have got: %v", *s.user, *res)

}
