package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yaruz/app/internal/domain/session"
	"strconv"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/minipkg/db/redis"
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

func (r SessionRepository) Key(userId uint) string {
	return fmt.Sprintf(keyFormatForSession, keyPrefixForSession, strconv.FormatUint(uint64(userId), 10))
}

// Get returns the Session with the specified user ID.
func (r SessionRepository) Get(ctx context.Context, userId uint) (*session.Session, error) {
	var entity session.Session
	res, err := r.db.DB().Get(ctx, r.Key(userId)).Result()

	if err != nil {
		if err == goredis.Nil {
			return nil, apperror.ErrNotFound
		}
		return nil, errors.Wrapf(apperror.ErrInternal, "Get() error: %v", err)
	}

	//err = entity.UnmarshalBinary([]byte(res))
	err = json.Unmarshal([]byte(res), &entity)
	if err != nil {
		return nil, errors.Wrapf(apperror.ErrInternal, "json.Unmarshal() error: %v", err)
	}

	return &entity, nil
}

// Create saves a new entity in the storage.
func (r SessionRepository) Create(ctx context.Context, entity *session.Session) error {
	return r.Set(ctx, entity)
}

func (r SessionRepository) Update(ctx context.Context, entity *session.Session) error {
	return r.Set(ctx, entity)
}

func (r SessionRepository) Set(ctx context.Context, entity *session.Session) error {
	//var _ encoding.BinaryMarshaler = entity

	if err := r.db.DB().Set(ctx, r.Key(entity.User.ID), entity, r.SessionLifeTime).Err(); err != nil {
		return errors.Wrapf(apperror.ErrInternal, "Create() error: %v", err)
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
