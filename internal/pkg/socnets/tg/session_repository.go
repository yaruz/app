package tg

import (
	"context"
	"errors"
	"github.com/yaruz/app/internal/pkg/apperror"

	"github.com/gotd/td/telegram"
)

type SessionRepositoryCreater func(sessionID string) SessionRepository

type SessionRepository interface {
	Get(ctx context.Context, ID string) (*Session, error)
	Set(ctx context.Context, entity *Session) error
	Delete(ctx context.Context, ID string) error
}

type SessionStorage struct {
	repository SessionRepository
	sessionID  string
}

var _ telegram.SessionStorage = (*SessionStorage)(nil)

func NewSessionStorage(repository SessionRepository, sessionID string) *SessionStorage {
	return &SessionStorage{
		repository: repository,
	}
}

func (s *SessionStorage) get(ctx context.Context) (*Session, error) {
	sess, err := s.repository.Get(ctx, s.sessionID)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return nil, err
		}
		sess = NewSession(s.sessionID)
	}
	return sess, nil
}

func (s *SessionStorage) LoadSession(ctx context.Context) ([]byte, error) {
	sess, err := s.get(ctx)
	if err != nil {
		return nil, err
	}
	return sess.Session, nil
}

func (s *SessionStorage) StoreSession(ctx context.Context, data []byte) error {
	sess, err := s.get(ctx)
	if err != nil {
		return err
	}
	sess.Session = data
	return s.repository.Set(ctx, sess)
}
