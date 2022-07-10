package tg

import (
	"context"
	"github.com/minipkg/log"

	mtproto_session "github.com/Kalinin-Andrey/mtproto/session"
	"github.com/Kalinin-Andrey/mtproto/telegram"

	"github.com/yaruz/app/internal/pkg/auth"
	"github.com/yaruz/app/internal/pkg/config"

	domain_session "github.com/yaruz/app/internal/domain/session"
	domain_tg_account "github.com/yaruz/app/internal/domain/tg_account"
)

const (
	warnChannelDefaultCapacity = 100
)

type IService interface {
	GetSessionLoader(session domain_session.Session) mtproto_session.SessionLoader
	IsAuthSessionRegistred(domainSession domain_session.Session) (bool, error)
	NewClient(domainSession domain_session.Session) (*telegram.Client, error)
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
	return s.domainSession.TgAccount.AuthSession, nil
}

func (s *AuthSessionLoader) Store(sess *mtproto_session.Session) error {
	var err error
	ctx := context.Background()

	if s.domainSession.TgAccount == nil {
		s.domainSession.TgAccount, err = s.tgAccountRepository.New(ctx)
		if err != nil {
			return err
		}
	}
	s.domainSession.TgAccount.AuthSession = sess

	if s.domainSession.TgAccount.ID == 0 {
		err = s.tgAccountRepository.Create(ctx, s.domainSession.TgAccount, s.domainSession.AccountSettings.LangID)
	} else {
		err = s.tgAccountRepository.Update(ctx, s.domainSession.TgAccount, s.domainSession.AccountSettings.LangID)
	}
	return err
}

type Service struct {
	logger              log.ILogger
	config              *config.Telegram
	authService         auth.Service
	sessionRepository   domain_session.Repository
	tgAccountRepository domain_tg_account.Repository
	Client              *telegram.Client
}

func NewService(logger log.ILogger, config *config.Telegram, authService auth.Service, sessionRepository domain_session.Repository, tgAccountRepository domain_tg_account.Repository) *Service {
	return &Service{
		logger:              logger,
		config:              config,
		authService:         authService,
		sessionRepository:   sessionRepository,
		tgAccountRepository: tgAccountRepository,
	}
}

func (s *Service) GetSessionLoader(domainSession domain_session.Session) mtproto_session.SessionLoader {
	return NewAuthSessionLoader(s.logger, s.authService, s.sessionRepository, s.tgAccountRepository, domainSession)
}

func (s *Service) IsAuthSessionRegistred(domainSession domain_session.Session) (bool, error) {
	if domainSession.TgAccount == nil || domainSession.TgAccount.ID == 0 || domainSession.TgAccount.AuthSession == nil {
		return false, nil
	}
	return true, nil
}

func (s *Service) NewClient(domainSession domain_session.Session) (*telegram.Client, error) {
	return telegram.NewClient(telegram.ClientConfig{
		SessionStorage:  s.GetSessionLoader(domainSession),
		ServerHost:      s.config.Application.ServerHost,
		PublicKeysFile:  s.config.Application.PublicKeyFile,
		AppID:           s.config.Application.AppID,
		AppHash:         s.config.Application.AppHash,
		InitWarnChannel: s.config.Application.InitWarnChannel,
	})
}
