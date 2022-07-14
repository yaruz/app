package tg

import (
	"context"
	"github.com/Kalinin-Andrey/mtproto"
	utils "github.com/Kalinin-Andrey/mtproto/examples/example_utils"
	"github.com/minipkg/db/redis"
	"github.com/minipkg/db/redis/cache"
	"github.com/minipkg/log"
	"github.com/pkg/errors"
	"github.com/xelaj/errs"
	"github.com/yaruz/app/internal/pkg/apperror"
	"time"

	mtproto_session "github.com/Kalinin-Andrey/mtproto/session"
	"github.com/Kalinin-Andrey/mtproto/telegram"

	"github.com/yaruz/app/internal/pkg/auth"
	"github.com/yaruz/app/internal/pkg/config"

	domain_session "github.com/yaruz/app/internal/domain/session"
	domain_tg_account "github.com/yaruz/app/internal/domain/tg_account"
)

const (
	warnChannelDefaultCapacity = 100
	ttlSentCode                = time.Hour
)

type IService interface {
	GetSessionLoader(session *domain_session.Session) mtproto_session.SessionLoader
	IsAuthSessionRegistred(domainSession *domain_session.Session) (bool, error)
	NewClient(domainSession *domain_session.Session) (*telegram.Client, error)
	AuthSendCode(ctx context.Context, client *telegram.Client, phone string) error
	AuthSignIn(ctx context.Context, client *telegram.Client, phone string, code string) error
	AuthCheckPassword(ctx context.Context, client *telegram.Client, password string) error
}

type AuthSessionLoader struct {
	logger              log.ILogger
	domainSession       *domain_session.Session
	authService         auth.Service
	sessionRepository   domain_session.Repository
	tgAccountRepository domain_tg_account.Repository
}

var _ mtproto_session.SessionLoader = (*AuthSessionLoader)(nil)

func NewAuthSessionLoader(logger log.ILogger, authService auth.Service, sessionRepository domain_session.Repository, tgAccountRepository domain_tg_account.Repository, domainSession *domain_session.Session) *AuthSessionLoader {
	return &AuthSessionLoader{
		logger:              logger,
		authService:         authService,
		sessionRepository:   sessionRepository,
		tgAccountRepository: tgAccountRepository,
		domainSession:       domainSession,
	}
}

func (s *AuthSessionLoader) Load() (*mtproto_session.Session, error) {
	if s.domainSession.TgAccount == nil || s.domainSession.TgAccount.AuthSession == nil {
		return nil, &errs.NotFoundError{}
	}
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
	redis               redis.IDB
	authService         auth.Service
	sessionRepository   domain_session.Repository
	tgAccountRepository domain_tg_account.Repository
	Client              *telegram.Client
}

var _ IService = (*Service)(nil)

func NewService(logger log.ILogger, config *config.Telegram, redis redis.IDB, authService auth.Service, sessionRepository domain_session.Repository, tgAccountRepository domain_tg_account.Repository) *Service {
	return &Service{
		logger:              logger,
		config:              config,
		redis:               redis,
		authService:         authService,
		sessionRepository:   sessionRepository,
		tgAccountRepository: tgAccountRepository,
	}
}

func (s *Service) GetSessionLoader(domainSession *domain_session.Session) mtproto_session.SessionLoader {
	return NewAuthSessionLoader(s.logger, s.authService, s.sessionRepository, s.tgAccountRepository, domainSession)
}

func (s *Service) IsAuthSessionRegistred(domainSession *domain_session.Session) (bool, error) {
	if domainSession.TgAccount == nil || domainSession.TgAccount.ID == 0 || domainSession.TgAccount.AuthSession == nil {
		return false, nil
	}
	return true, nil
}

func (s *Service) NewClient(domainSession *domain_session.Session) (*telegram.Client, error) {
	client, err := telegram.NewClient(telegram.ClientConfig{
		SessionStorage:  s.GetSessionLoader(domainSession),
		ServerHost:      s.config.Application.ServerHost,
		PublicKeysFile:  s.config.Application.PublicKeyFile,
		AppID:           s.config.Application.AppID,
		AppHash:         s.config.Application.AppHash,
		InitWarnChannel: s.config.Application.InitWarnChannel,
	})
	if err != nil {
		return nil, err
	}
	client.Warnings = make(chan error)
	utils.ReadWarningsToStdErr(client.Warnings) // todo: придумать что-то получше
	return client, nil
}

func (s *Service) AuthSendCode(ctx context.Context, client *telegram.Client, phone string) error {
	sentCode, err := client.AuthSendCode(
		phone, int32(s.config.Application.AppID), s.config.Application.AppHash, &telegram.CodeSettings{},
	)

	// this part shows how to deal with errors (if you want of course. No one
	// like errors, but the can be return sometimes)
	if err != nil {
		errResponse := &mtproto.ErrResponseCode{}
		if !errors.As(err, &errResponse) {
			// some strange error, looks like a bug actually
			return errors.Wrapf(apperror.ErrInternal, "AuthSendCode got an unexpected error")
		} else {
			/*if errResponse.Message == "AUTH_RESTART" {
				println("Oh crap! You accidentally restart authorization process!")
				println("You should login only once, if you'll spam 'AuthSendCode' method, you can be")
				println("timeouted to loooooooong long time. You warned.")
			} else if errResponse.Message == "FLOOD_WAIT_X" {
				println("No way... You've reached flood timeout! Did i warn you? Yes, i am. That's what")
				println("happens, when you don't listen to me...")
				println()
				timeoutDuration := time.Second * time.Duration(errResponse.AdditionalInfo.(int))

				println("Repeat after " + timeoutDuration.String())
			} else {
				println("Oh crap! Got strange error:")
				pp.Println(errResponse)
			}*/

			return errors.Wrapf(apperror.ErrBadRequest, errResponse.Message)
		}
	}

	return s.redis.CacheSet(&cache.Item{
		Ctx:   ctx,
		Key:   phone + "_sentCode",
		Value: sentCode,
		TTL:   ttlSentCode,
	})
}

func (s *Service) AuthSignIn(ctx context.Context, client *telegram.Client, phone string, code string) error {
	sentCode := &telegram.AuthSentCode{}
	if err := s.redis.CacheGet(ctx, phone+"_sentCode", sentCode); err != nil {
		return err
	}

	auth, err := client.AuthSignIn(
		phone,
		sentCode.PhoneCodeHash,
		code,
	)
	s.logger.Debug(auth)
	if err == nil {
		return nil
	}

	// could be some errors:
	errResponse := &mtproto.ErrResponseCode{}
	ok := errors.As(err, &errResponse)
	// checking that error type is correct, and error msg is actualy ask for password
	if !ok || errResponse.Message != "SESSION_PASSWORD_NEEDED" {
		s.logger.Errorf("tg.Service.AuthSignIn() error: %v; auth: %v", err, auth)
		return err
	}

	return ErrSessionPasswordNeeded
}

func (s *Service) AuthCheckPassword(ctx context.Context, client *telegram.Client, password string) error {
	accountPassword, err := client.AccountGetPassword()
	if err != nil {
		return err
	}

	// GetInputCheckPassword is fast response object generator
	inputCheck, err := telegram.GetInputCheckPassword(password, accountPassword)
	if err != nil {
		return err
	}

	auth, err := client.AuthCheckPassword(inputCheck)
	if err != nil {
		return err
	}
	s.logger.Debug(auth)

	return nil
}
