package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/minipkg/db/redis"

	session_proto "github.com/yaruz/app/internal/app/proto/session"
	"github.com/yaruz/app/internal/domain/session"
	"github.com/yaruz/app/internal/pkg/apperror"
)

const (
	keyFormatForSession = "%s%s"
	keyPrefixForSession = "session_"
)

// SessionRepository is a repository for the session entity
type SessionRepository struct {
	repository
	SessionLifeTime time.Duration
}

var _ session.Repository = (*SessionRepository)(nil)

// NewSessionRepository creates a new SessionRepository
func NewSessionRepository(dbase redis.IDB, sessionLifeTimeInHours uint) (*SessionRepository, error) {
	r := &SessionRepository{
		repository: repository{
			db: dbase,
		},
		SessionLifeTime: time.Duration(int64(sessionLifeTimeInHours)) * time.Hour,
	}
	return r, nil
}

func (r *SessionRepository) Key(userID uint) string {
	return fmt.Sprintf(keyFormatForSession, keyPrefixForSession, strconv.FormatUint(uint64(userID), 10))
}

// Get returns the Session with the specified user ID.
func (r *SessionRepository) Get(ctx context.Context, userId uint) (*session.Session, error) {
	//var entity session.Session
	res, err := r.db.DB().Get(ctx, r.Key(userId)).Result()

	if err != nil {
		if err == goredis.Nil {
			return nil, apperror.ErrNotFound
		}
		return nil, errors.Wrapf(apperror.ErrInternal, "Get() error: %v", err)
	}

	//err = entity.UnmarshalBinary([]byte(res))
	entity, err := session_proto.UnmarshalBinary([]byte(res))
	if err != nil {
		return nil, errors.Wrapf(apperror.ErrInternal, "session_proto.UnmarshalBinary() error: %v", err)
	}

	return entity, nil
}

// Create saves a new entity in the storage.
func (r *SessionRepository) Create(ctx context.Context, entity *session.Session) error {
	return r.Set(ctx, entity)
}

func (r *SessionRepository) Update(ctx context.Context, entity *session.Session) error {
	return r.Set(ctx, entity)
}

func (r *SessionRepository) Set(ctx context.Context, entity *session.Session) error {
	//var in interface{} = entity
	//v, ok := in.(encoding.BinaryMarshaler)
	//if !ok {
	//	return errors.Errorf("Can not cast entity session %v to the interface encoding.BinaryMarshaler", v)
	//}
	//
	//if err := r.db.DB().Set(ctx, r.Key(entity.User.ID), v, r.SessionLifeTime).Err(); err != nil {
	data, err := session_proto.MarshalBinary(entity)
	if err != nil {
		return errors.Wrapf(apperror.ErrInternal, "session_proto.MarshalBinary() error: %v", err)
	}

	if err := r.db.DB().Set(ctx, r.Key(entity.User.ID), data, r.SessionLifeTime).Err(); err != nil {
		return errors.Wrapf(apperror.ErrInternal, "Redis Set() error: %v", err)
	}
	return nil
}

// Delete removes the entity with given ID from the storage.
func (r *SessionRepository) Delete(ctx context.Context, entity *session.Session) error {

	if err := r.db.DB().Del(ctx, r.Key(entity.User.ID)).Err(); err != nil {
		return errors.Wrapf(apperror.ErrInternal, "Delete error: %v", err)
	}
	return nil
}
