package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"

	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/internal/pkg/socnets/tg"

	tg_session_proto "github.com/yaruz/app/internal/app/proto/tg_session"
)

const (
	tgSession_KeyPrefix = "tg_session_"
)

// SessionRepository is a repository for the tg session
type TgSessionRepository struct {
	Repository
	ID string
}

var _ tg.SessionRepository = (*TgSessionRepository)(nil)

// NewTgSessionRepository creates a new TgSessionRepository
func NewTgSessionRepository(repository Repository) *TgSessionRepository {
	return &TgSessionRepository{
		Repository: repository,
	}
}

func (r *TgSessionRepository) key(ID string) string {
	return tgSession_KeyPrefix + ID
}

// Get returns the TgSession with the specified ID.
func (r *TgSessionRepository) Get(ctx context.Context, ID string) (*tg.Session, error) {
	res, err := r.DB().Get(ctx, r.key(ID)).Result()

	if err != nil {
		if err == goredis.Nil {
			return nil, apperror.ErrNotFound
		}
		return nil, fmt.Errorf("[%w] Get() error: %s", apperror.ErrInternal, err.Error())
	}

	entity, err := tg_session_proto.UnmarshalBinary([]byte(res))
	if err != nil {
		return nil, fmt.Errorf("[%w] tg_session_proto.UnmarshalBinary() error: %s", apperror.ErrInternal, err.Error())
	}

	return entity, nil
}

func (r *TgSessionRepository) Set(ctx context.Context, entity *tg.Session) error {
	data, err := tg_session_proto.MarshalBinary(entity)
	if err != nil {
		return fmt.Errorf("[%w] tg_session_proto.MarshalBinary() error: %s", apperror.ErrInternal, err.Error())
	}

	if err := r.DB().Set(ctx, r.key(entity.ID), data, 0).Err(); err != nil {
		return fmt.Errorf("[%w] Redis Set() error: %s", apperror.ErrInternal, err.Error())
	}
	return nil
}

// Delete removes the entity with given ID from the storage.
func (r *TgSessionRepository) Delete(ctx context.Context, ID string) error {

	if err := r.DB().Del(ctx, r.key(ID)).Err(); err != nil {
		return fmt.Errorf("[%w] Delete error: %s", apperror.ErrInternal, err.Error())
	}
	return nil
}
