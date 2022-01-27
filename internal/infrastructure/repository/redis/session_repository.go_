package redis

import (
	"context"
	"encoding"
	"fmt"
	"strconv"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/internal/pkg/auth"
	"github.com/yaruz/app/internal/pkg/session"

	"github.com/minipkg/db/redis"

	"github.com/yaruz/app/internal/domain/user"
)

const (
	keyPrefixForSession = "session_"
)

// SessionRepository is a repository for the session entity
type SessionRepository struct {
	repository
	UserRepo        user.Repository
	SessionLifeTime time.Duration
}

var _ auth.SessionRepository = (*SessionRepository)(nil)

// New creates a new SessionRepository
func NewSessionRepository(dbase redis.IDB, sessionLifeTimeInHours uint, userRepo user.Repository) (*SessionRepository, error) {
	r := &SessionRepository{
		repository: repository{
			db: dbase,
		},
		UserRepo:        userRepo,
		SessionLifeTime: time.Duration(int64(sessionLifeTimeInHours)) * time.Hour,
	}
	return r, nil
}

func (r SessionRepository) Key(userId uint) string {
	return fmt.Sprintf("%s%s", keyPrefixForSession, strconv.FormatUint(uint64(userId), 10))
}

func (r SessionRepository) NewEntity(ctx context.Context, userId uint) (*session.Session, error) {
	user, err := r.UserRepo.Get(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &session.Session{
		UserID: userId,
		User:   *user,
	}, nil
}

func (r SessionRepository) GetData(session *session.Session) session.Data {
	return session.Data
}

func (r *SessionRepository) SetData(session *session.Session, data session.Data) error {
	session.Data = data
	return r.Save(session)
}

func (r *SessionRepository) Save(session *session.Session) error {
	return r.Update(session.Ctx, session)
}

// GetByUserID returns the Session with the specified user ID.
func (r SessionRepository) Get(ctx context.Context, userId uint) (*session.Session, error) {
	var entity session.Session
	res, err := r.db.DB().Get(ctx, r.Key(userId)).Result()

	if err != nil {
		if err == goredis.Nil {
			return nil, apperror.ErrNotFound
		}
		return nil, errors.Wrapf(apperror.ErrInternal, "Get() error: %v", err)
	}

	err = entity.UnmarshalBinary([]byte(res))
	if err != nil {
		return nil, errors.Wrapf(apperror.ErrInternal, "json.Unmarshal() error: %v", err)
	}

	return &entity, nil
}

// Create saves a new entity in the storage.
func (r SessionRepository) Create(ctx context.Context, entity *session.Session) error {
	var _ encoding.BinaryMarshaler = entity

	if err := r.db.DB().Set(ctx, r.Key(entity.User.ID), entity, r.SessionLifeTime).Err(); err != nil {
		return errors.Wrapf(apperror.ErrInternal, "Create() error: %v", err)
	}
	return nil
}

// Update updates the entity with given ID in the storage.
func (r SessionRepository) Update(ctx context.Context, entity *session.Session) error {

	if err := r.db.DB().Set(ctx, r.Key(entity.User.ID), entity, r.SessionLifeTime).Err(); err != nil {
		return errors.Wrapf(apperror.ErrInternal, "Update() error: %v", err)
	}
	return nil
}

// Delete removes the entity with given ID from the storage.
func (r SessionRepository) Delete(ctx context.Context, entity *session.Session) error {

	if err := r.db.DB().Del(ctx, r.Key(entity.User.ID)).Err(); err != nil {
		return errors.Wrapf(apperror.ErrInternal, "Delete error: %v", err)
	}
	return nil
}
