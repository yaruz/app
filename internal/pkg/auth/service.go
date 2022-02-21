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
	SignIn(ctx context.Context, code, state string, langId uint) (context.Context, error)
	// authenticate authenticates a user using username and password.
	// It returns a JWT token if authentication succeeds. Otherwise, an error is returned.
	//Login(ctx context.Context, username, password string) (string, error)
	//Register(ctx context.Context, username, password string) (string, error)
	//NewUser(username, password string) (*user.User, error)
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

func (s service) NewSession(ctx context.Context, jwtClaims *auth.Claims, user *user.User, langId uint) (*session.Session, error) {
	if user == nil {
		var err error

		if user, err = s.userService.GetByAccountID(ctx, jwtClaims.User.Id, langId); err != nil {
			return nil, err
		}
	}

	return &session.Session{
		User:      user,
		JwtClaims: jwtClaims,
	}, nil
}

func (s service) createSession(ctx context.Context, jwtClaims *auth.Claims, user *user.User, langId uint) (context.Context, *session.Session, error) {
	sess, err := s.NewSession(ctx, jwtClaims, user, langId)
	if err != nil {
		return ctx, nil, err
	}

	err = s.session.Create(ctx, sess)
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

func (s service) updateSession(ctx context.Context, sess *session.Session, jwtClaims *auth.Claims, user *user.User, langId uint) (context.Context, *session.Session, error) {
	if user != nil {
		sess.User = user
	}
	sess.JwtClaims = jwtClaims

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

func (s service) accountSetUserID(user *auth.User, userID uint) {
	user.Properties[s.Application+"ID"] = strconv.Itoa(int(userID))
}

func (s service) accountGetUserID(user *auth.User) (uint, error) {
	var userID uint

	if id, ok := user.Properties[s.Application+"ID"]; !ok {
		return 0, apperror.ErrNotFound
		ID, err := strconv.Atoi(id)
		if err != nil {
			return 0, err
		}
		userID = uint(ID)
	}
	return userID, nil
}

func (s service) SessionInit(ctx context.Context, token string, langID uint) (context.Context, error) {
	//todo: использование нахождения юзера по AccountID
	//todo: настройки пользователя и работа с langID
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
		user, err := s.userService.Get(ctx, userID, langID)
		if err != nil {
			if !errors.Is(err, apperror.ErrNotFound) {
				return ctx, err
			}
			return ctx, errors.Wrapf(apperror.ErrNotFound, "User with ID = %d not found", userID)
		}

		ctx, _, err = s.createSession(ctx, jwtClaims, user, langID)
		if err != nil {
			return ctx, err
		}
	} else {
		ctx, _, err = s.updateSession(ctx, sess, jwtClaims, nil, langID)
	}

	return ctx, nil
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) SignIn(ctx context.Context, code, state string, langId uint) (context.Context, error) {
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

	user, err := s.userService.GetByAccountID(ctx, jwtClaims.User.Id, langId)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return ctx, err
		}
		user, err = s.signUp(ctx, jwtClaims, langId)
	}

	if userID, err := s.accountGetUserID(&jwtClaims.User); err == nil {
		if userID != user.ID {
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

	ctx, err = s.SessionInit(ctx, token.AccessToken, langId)

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
