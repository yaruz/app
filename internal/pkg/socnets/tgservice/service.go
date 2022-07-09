package tgservice

import (
	mtproto_session "github.com/Kalinin-Andrey/mtproto/session"
	"github.com/minipkg/log"
	"github.com/yaruz/app/internal/pkg/auth"

	domain_session "github.com/yaruz/app/internal/domain/session"
	domain_tg_account "github.com/yaruz/app/internal/domain/tg_account"
)

type IService interface {
	GetSessionLoader(session domain_session.Session) mtproto_session.SessionLoader
}

type AuthSessionLoader struct {
	logger              log.ILogger
	domainSession       domain_session.Session
	authService         auth.Service
	sessionRepository   domain_session.Repository
	tgAccountRepository domain_tg_account.Repository
}

var _ mtproto_session.SessionLoader = (*AuthSessionLoader)(nil)

func NewAuthSessionLoader(logger log.ILogger, authService auth.Service, sessionRepository domain_session.Repository, tgAccountRepository domain_tg_account.Repository, domainSession domain_session.Session) *AuthSessionLoader {
	return &AuthSessionLoader{
		logger:              logger,
		authService:         authService,
		sessionRepository:   sessionRepository,
		tgAccountRepository: tgAccountRepository,
		domainSession:       domainSession,
	}
}

func (s *AuthSessionLoader) Load() (*mtproto_session.Session, error) {
	return nil, nil
}

func (s *AuthSessionLoader) Store(sess *mtproto_session.Session) error {
	return nil
}

type Service struct {
	logger              log.ILogger
	authService         auth.Service
	sessionRepository   domain_session.Repository
	tgAccountRepository domain_tg_account.Repository
}

func NewService(logger log.ILogger, authService auth.Service, sessionRepository domain_session.Repository, tgAccountRepository domain_tg_account.Repository) *Service {
	return &Service{
		logger:              logger,
		authService:         authService,
		sessionRepository:   sessionRepository,
		tgAccountRepository: tgAccountRepository,
	}
}

func (s *Service) GetSessionLoader(domainSession domain_session.Session) mtproto_session.SessionLoader {
	return NewAuthSessionLoader(s.logger, s.authService, s.sessionRepository, s.tgAccountRepository, domainSession)
}
