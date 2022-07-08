package tgservice

import (
	"github.com/Kalinin-Andrey/mtproto/session"
	"github.com/minipkg/log"

	domain_session "github.com/yaruz/app/internal/domain/session"
	domain_tg_account "github.com/yaruz/app/internal/domain/tg_account"
)

type IService interface {
	GetSessionLoader() session.SessionLoader
}

type SessionLoader struct {
	sessionRepository   domain_session.Repository
	tgAccountRepository domain_tg_account.Repository
}

var _ session.SessionLoader = (*SessionLoader)(nil)

func NewSessionLoader(sessionRepository domain_session.Repository, tgAccountRepository domain_tg_account.Repository) *SessionLoader {
	return &SessionLoader{
		sessionRepository:   sessionRepository,
		tgAccountRepository: tgAccountRepository,
	}
}

func (s *SessionLoader) Load() (*session.Session, error) {
	return nil, nil
}

func (s *SessionLoader) Store(sess *session.Session) error {
	return nil
}

type Service struct {
	logger              log.ILogger
	sessionRepository   domain_session.Repository
	tgAccountRepository domain_tg_account.Repository
}

func NewService(logger log.ILogger, sessionRepository domain_session.Repository, tgAccountRepository domain_tg_account.Repository) *Service {
	return &Service{
		logger:              logger,
		sessionRepository:   sessionRepository,
		tgAccountRepository: tgAccountRepository,
	}
}

func (s *Service) GetSessionLoader() session.SessionLoader {
	return NewSessionLoader(s.sessionRepository, s.tgAccountRepository)
}
