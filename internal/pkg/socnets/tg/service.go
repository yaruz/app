package tg

import (
	"context"
	"errors"
	"fmt"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/minipkg/log"
)

type Service interface {
	IsAuth(ctx context.Context, sessionID string) (bool, error)
	SendCode(ctx context.Context, sessionID string, phone string) error
	SignIn(ctx context.Context, sessionID string, code string, password string) (*tg.User, error)
}

type Config struct {
	Application Application
}

type Application struct {
	Title           string
	ShortName       string
	AppID           int
	AppHash         string
	ServerHost      string
	PublicKeyFile   string
	InitWarnChannel bool
}

type service struct {
	config            *Config
	logger            log.Logger
	sessionRepository SessionRepository
}

var _ Service = (*service)(nil)

func New(config *Config, logger log.Logger, sessionRepository SessionRepository) *service {
	return &service{
		config:            config,
		logger:            logger,
		sessionRepository: sessionRepository,
	}
}

func (s *service) newClient(sessionID string) *telegram.Client {
	return telegram.NewClient(
		s.config.Application.AppID,
		s.config.Application.AppHash,
		telegram.Options{
			SessionStorage: NewSessionStorage(s.sessionRepository, sessionID),
			Logger:         s.logger.ZapLogger(),
		},
	)
}

// IsAuth checks auth status.
func (s *service) IsAuth(ctx context.Context, sessionID string) (bool, error) {
	var status *auth.Status
	var err error
	client := s.newClient(sessionID)

	err = client.Run(ctx, func(ctx context.Context) error {
		status, err = client.Auth().Status(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("[%w] get tg Status error: %s", ErrInternal, err.Error())
	}
	return status.Authorized, nil
}

func (s *service) SendCode(ctx context.Context, sessionID string, phone string) error {
	client := s.newClient(sessionID)
	sess, err := s.getSession(ctx, sessionID)
	if err != nil {
		return err
	}
	sess.Phone = phone

	return client.Run(ctx, func(ctx context.Context) error {
		sentCode, err := client.Auth().SendCode(ctx, phone, auth.SendCodeOptions{})
		if err != nil {
			return fmt.Errorf("[%w] SendCode error: %s", ErrInternal, err.Error())
		}

		switch sc := sentCode.(type) {
		case *tg.AuthSentCode:
			sess.PhoneCodeHash = sc.PhoneCodeHash
			if err = s.sessionRepository.Set(ctx, sess); err != nil {
				return fmt.Errorf("[%w] SessionRepository.Set error: %s", ErrInternal, err.Error())
			}
		case *tg.AuthSentCodeSuccess:
			switch a := sc.Authorization.(type) {
			case *tg.AuthAuthorization:
				// Looks that we are already authorized.
				return nil
			case *tg.AuthAuthorizationSignUpRequired:
				// SignUpRequired error
				return ErrSignUpRequired
			default:
				return fmt.Errorf("[%w] Unexpected authorization type: %T", ErrInternal, a)
			}
		default:
			return fmt.Errorf("[%w] Unexpected sent code type: %T", ErrInternal, sentCode)
		}
		return nil
	})
}

func (s *service) SignIn(ctx context.Context, sessionID string, code string, password string) (*tg.User, error) {
	var user *tg.User
	client := s.newClient(sessionID)
	sess, err := s.getSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if sess.Phone == "" {
		return nil, fmt.Errorf("[%w] Tg session phone is empty", ErrBadRequest)
	}
	if sess.PhoneCodeHash == "" {
		return nil, fmt.Errorf("[%w] Tg session phone code hash is empty", ErrBadRequest)
	}

	err = client.Run(ctx, func(ctx context.Context) error {
		tgAuth, err := client.Auth().SignIn(ctx, sess.Phone, code, sess.PhoneCodeHash)
		if err != nil {
			if !errors.Is(err, auth.ErrPasswordAuthNeeded) {
				return fmt.Errorf("[%w] Tg SignIn error: %s", ErrBadParams, err.Error())
			}

			if password == "" {
				// ошибка ErrPasswordAuthNeeded
				return ErrPasswordAuthNeeded
			}

			if tgAuth, err = client.Auth().Password(ctx, password); err != nil {
				return fmt.Errorf("[%w] Tg SignIn with password error: %s", ErrBadParams, err.Error())
			}
		}

		u := tgAuth.User
		switch obj := u.(type) {
		case *tg.UserEmpty:
			// ошибка UserEmpty
			return ErrUserEmpty
		case *tg.User:
			user = obj
			return nil
		default:
			return fmt.Errorf("[%w] type casting tg user error: %s", ErrInternal, err.Error())
		}
		return nil
	})
	return user, err
}

func (s *service) getSession(ctx context.Context, sessionID string) (*Session, error) {
	sess, err := s.sessionRepository.Get(ctx, sessionID)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return nil, err
		}
		sess = NewSession(sessionID)
	}
	return sess, nil
}
