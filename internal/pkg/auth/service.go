package auth

import (
	"context"
	_ "embed"
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/casdoor/casdoor-go-sdk/auth"
	"github.com/minipkg/log"
	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/internal/pkg/config"

	"github.com/yaruz/app/internal/domain/session"
	"github.com/yaruz/app/internal/domain/user"
)

// Service encapsulates the authentication logic.
type Service interface {
	GetSignUpUrl() string
	GetSignInUrl() string
	GetForgetUrl() string
	GetSession(ctx context.Context) *session.Session
	UpdateSession(ctx context.Context, sess *session.Session) (context.Context, *session.Session, error)
	AccountSettingsUpdate(ctx context.Context, accountSettings *user.AccountSettings) (context.Context, error)
	SessionInit(ctx context.Context, token string, accountSettings *user.AccountSettings) (context.Context, error)
	SignIn(ctx context.Context, code, state string, accountSettings *user.AccountSettings) (context.Context, error)
	StringTokenValidation(stringToken string) error
}

var _ Service = service{}

type UserService interface {
}

type service struct {
	//signingKey        string
	//tokenExpiration   uint
	userService user.IService
	logger      log.ILogger
	session     session.Repository
	//tokenRepository   TokenRepository
	Endpoint          string
	ClientId          string
	ClientSecret      string
	Organization      string
	Application       string
	SignInRedirectURL string
	JWTSigningKey     string
	JWTExpiration     uint
	SessionlifeTime   uint
}

type contextKey int

const (
	userSessionKey contextKey = iota
)

//go:embed token_jwt_key.pem
var JwtPublicKey string

// NewService creates a new authentication service.
func NewService(logger log.ILogger, cfg config.Auth, userService user.IService, session session.Repository) *service {
	s := &service{
		logger:            logger,
		Endpoint:          cfg.Endpoint,
		ClientId:          cfg.ClientId,
		ClientSecret:      cfg.ClientSecret,
		Organization:      cfg.Organization,
		Application:       cfg.Application,
		SignInRedirectURL: cfg.SignInRedirectURL,
		JWTSigningKey:     cfg.JWTSigningKey,
		JWTExpiration:     cfg.JWTExpiration,
		SessionlifeTime:   cfg.SessionlifeTime,
		userService:       userService,
		session:           session,
	}
	auth.InitConfig(s.Endpoint, s.ClientId, s.ClientSecret, JwtPublicKey, s.Organization, s.Application)
	return s
}

func (s service) newSession(ctx context.Context, jwtClaims *auth.Claims, user *user.User, accountSettings *user.AccountSettings) *session.Session {
	return &session.Session{
		User:            user,
		JwtClaims:       jwtClaims,
		AccountSettings: accountSettings,
	}
}

func (s service) createSession(ctx context.Context, jwtClaims *auth.Claims, user *user.User, accountSettings *user.AccountSettings) (context.Context, *session.Session, error) {
	if jwtClaims == nil || user == nil || accountSettings == nil {
		return ctx, nil, errors.Wrapf(apperror.ErrBadParams, "jwtClaims == %v \nuser == %v \naccountSettings == %v", jwtClaims, user, accountSettings)
	}
	sess := s.newSession(ctx, jwtClaims, user, accountSettings)

	err := s.session.Create(ctx, sess)
	if err != nil {
		return ctx, nil, err
	}

	ctx = context.WithValue(
		ctx,
		userSessionKey,
		sess,
	)

	return ctx, sess, nil
}

func (s service) UpdateSession(ctx context.Context, sess *session.Session) (context.Context, *session.Session, error) {
	err := s.session.Update(ctx, sess)
	if err != nil {
		return ctx, nil, err
	}

	ctx = context.WithValue(
		ctx,
		userSessionKey,
		sess,
	)

	return ctx, sess, nil
}

func (s service) updateSession(ctx context.Context, sess *session.Session, jwtClaims *auth.Claims, user *user.User, accountSettings *user.AccountSettings) (context.Context, *session.Session, error) {
	if user != nil {
		sess.User = user
	}

	if accountSettings != nil {
		sess.AccountSettings = accountSettings
	}

	if jwtClaims != nil {
		sess.JwtClaims = jwtClaims
	}

	return s.UpdateSession(ctx, sess)
}

func (s service) GetSession(ctx context.Context) *session.Session {
	return ctx.Value(userSessionKey).(*session.Session)
}

func (s service) GetSignUpUrl() string {
	return fmt.Sprintf("%s/signup/%s", s.Endpoint, s.Application)
}

func (s service) GetSignInUrl() string {
	return fmt.Sprintf("%s/login/oauth/authorize?client_id=%s&response_type=code&scope=read&state=%s&redirect_uri=%s", s.Endpoint, s.ClientId, s.Application, s.SignInRedirectURL)
}

func (s service) GetForgetUrl() string {
	return fmt.Sprintf("%s/forget/%s", s.Endpoint, s.Application)
}

func (s service) accountGetUintParam(account *auth.User, paramName string) (uint, error) {
	var param uint

	if pStr, ok := account.Properties[paramName]; !ok {
		return 0, apperror.ErrNotFound
		pInt, err := strconv.Atoi(pStr)
		if err != nil {
			return 0, err
		}
		param = uint(pInt)
	}
	return param, nil
}

func (s service) accountSetUserID(account *auth.User, userID uint) {
	account.Properties[s.Application+"ID"] = strconv.Itoa(int(userID))
}

func (s service) accountGetUserID(account *auth.User) (uint, error) {
	userID, err := s.accountGetUintParam(account, s.Application+"ID")
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s service) accountSetSettings(account *auth.User, accountSettings *user.AccountSettings) {
	account.Properties["langID"] = strconv.Itoa(int(accountSettings.LangID))
}

func (s service) accountGetSettings(account *auth.User) (*user.AccountSettings, error) {
	langID, err := s.accountGetUintParam(account, "langID")
	if err != nil {
		return nil, err
	}

	return &user.AccountSettings{
		LangID: langID,
	}, nil
}

func (s service) AccountSettingsUpdate(ctx context.Context, accountSettings *user.AccountSettings) (context.Context, error) {
	sess := s.GetSession(ctx)

	if sess.AccountSettings.LangID == accountSettings.LangID {
		return ctx, nil
	}

	user, err := s.userService.Get(ctx, sess.User.ID, accountSettings.LangID)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return ctx, err
		}
		return ctx, errors.Wrapf(apperror.ErrNotFound, "User with ID = %d not found", sess.User.ID)
	}

	s.accountSetSettings(&sess.JwtClaims.User, accountSettings)

	ctx, _, err = s.updateSession(ctx, sess, nil, user, accountSettings)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (s service) SessionInit(ctx context.Context, token string, accountSettings *user.AccountSettings) (context.Context, error) {
	// todo: +передать токен на фронт
	jwtClaims, err := auth.ParseJwtToken(token)
	if err != nil {
		return ctx, err
	}
	jwtClaims.AccessToken = token

	userID, err := s.accountGetUserID(&jwtClaims.User)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return ctx, err
		}
		return ctx, errors.Wrapf(apperror.ErrNotFound, "Param %q not found in JwtClaims.User = %v", s.Application+"ID", jwtClaims.User)
	}

	sess, err := s.session.Get(ctx, userID)
	if err != nil && !errors.Is(err, apperror.ErrNotFound) {
		return ctx, err
	}

	if sess == nil {
		user, err := s.userService.Get(ctx, userID, accountSettings.LangID)
		if err != nil {
			if !errors.Is(err, apperror.ErrNotFound) {
				return ctx, err
			}
			return ctx, errors.Wrapf(apperror.ErrNotFound, "User with ID = %d not found", userID)
		}

		ctx, _, err = s.createSession(ctx, jwtClaims, user, accountSettings)
		if err != nil {
			return ctx, err
		}
	} else {
		ctx, _, err = s.updateSession(ctx, sess, jwtClaims, nil, nil)
		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) SignIn(ctx context.Context, code, state string, accountSettings *user.AccountSettings) (context.Context, error) {
	token, err := auth.GetOAuthToken(code, state)
	if err != nil {
		return ctx, err
	}

	if err = s.StringTokenValidation(token.AccessToken); err != nil {
		return ctx, err
	}

	jwtClaims, err := auth.ParseJwtToken(token.AccessToken)
	if err != nil {
		return ctx, err
	}
	jwtClaims.AccessToken = token.AccessToken

	user, err := s.userService.GetByAccountID(ctx, jwtClaims.User.Id, accountSettings.LangID)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return ctx, err
		}
		user, err = s.signUp(ctx, jwtClaims, accountSettings.LangID)
	}

	if userID, err := s.accountGetUserID(&jwtClaims.User); err == nil {
		if userID != user.ID {
			//	user был создан, но был удалён без удаления аккаунта
			return ctx, err
		}
	} else if errors.Is(err, apperror.ErrNotFound) {
		s.accountSetUserID(&jwtClaims.User, user.ID)
		ok, err := auth.UpdateUserForColumns(&jwtClaims.User, []string{"properties"})
		if !ok || err != nil {
			return ctx, errors.Wrapf(apperror.ErrInternal, fmt.Sprintf("User not updated. Ok = %t, err = %q.", ok, err.Error()))
		}
	} else {
		if err != nil {
			return ctx, err
		}
	}

	ctx, err = s.SessionInit(ctx, token.AccessToken, accountSettings)

	//affected, err := object.UpdateMemberOnlineStatus(&claims.User, true, util.GetCurrentTime())
	//if err != nil {
	//	c.ResponseError(err.Error())
	//	return
	//}

	return ctx, nil
}

func (s service) signUp(ctx context.Context, jwtClaims *auth.Claims, langId uint) (*user.User, error) {
	user, err := s.userService.New(ctx)
	if err != nil {
		return nil, err
	}

	err = user.SetAccountID(ctx, jwtClaims.User.Id)
	if err != nil {
		return nil, err
	}

	err = user.SetEmail(ctx, jwtClaims.User.Email)
	if err != nil {
		return nil, err
	}

	return user, s.userService.Create(ctx, user, langId)
}

func (s service) StringTokenValidation(stringToken string) error {
	//	temporary validation method
	_, err := auth.ParseJwtToken(stringToken)
	if err != nil {
		return err
	}
	return nil
}
