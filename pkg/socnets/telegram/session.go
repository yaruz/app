package telegram

import (
	"github.com/Kalinin-Andrey/mtproto/session"
)

type Session struct {
	Auth *session.Session
}

type SessionLoader struct {
}

var _ session.SessionLoader = (*SessionLoader)(nil)

func (s *SessionLoader) Load() (*session.Session, error) {
	return nil, nil
}

func (s *SessionLoader) Store(sess *session.Session) error {
	return nil
}
