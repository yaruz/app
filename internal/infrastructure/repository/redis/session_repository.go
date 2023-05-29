package redis

import (
	"context"
	"time"

	"github.com/pkg/errors"
	goredis "github.com/redis/go-redis/v9"

	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/internal/pkg/session"

	session_proto "github.com/yaruz/app/internal/app/proto/session"
)

const (
	session_KeyPrefix = "session_"
)

// SessionRepository is a repository for the session entity
type SessionRepository struct {
	Repository
	SessionLifeTime time.Duration
}

var _ session.Repository = (*SessionRepository)(nil)

// NewSessionRepository creates a new SessionRepository
func NewSessionRepository(repository Repository, sessionLifeTimeInHours uint) (*SessionRepository, error) {
	r := &SessionRepository{
		Repository:      repository,
		SessionLifeTime: time.Duration(int64(sessionLifeTimeInHours)) * time.Hour,
	}
	return r, nil
}

func (r *SessionRepository) key(ID string) string {
	return session_KeyPrefix + ID
}

// Get returns the Session with the specified ID.
func (r *SessionRepository) Get(ctx context.Context, ID string) (*session.Session, error) {
	//var entity session.Session
	res, err := r.DB().Get(ctx, r.key(ID)).Result()

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

func (r *SessionRepository) Set(ctx context.Context, entity *session.Session) error {
	//var in interface{} = entity
	//v, ok := in.(encoding.BinaryMarshaler)
	//if !ok {
	//	return errors.Errorf("Can not cast entity session %v to the interface encoding.BinaryMarshaler", v)
	//}
	//
	//if err := r.db.DB().Set(ctx, r.key(entity.User.ID), v, r.SessionLifeTime).Err(); err != nil {
	data, err := session_proto.MarshalBinary(entity)
	if err != nil {
		return errors.Wrapf(apperror.ErrInternal, "session_proto.MarshalBinary() error: %v", err)
	}

	if err := r.DB().Set(ctx, r.key(entity.ID), data, r.SessionLifeTime).Err(); err != nil {
		return errors.Wrapf(apperror.ErrInternal, "Redis Set() error: %v", err)
	}
	return nil
}

// Delete removes the entity with given ID from the storage.
func (r *SessionRepository) Delete(ctx context.Context, ID string) error {

	if err := r.DB().Del(ctx, r.key(ID)).Err(); err != nil {
		return errors.Wrapf(apperror.ErrInternal, "Delete error: %v", err)
	}
	return nil
}
