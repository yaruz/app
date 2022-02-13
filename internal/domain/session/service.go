package session

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/casdoor/casdoor-go-sdk/auth"
	auth2 "github.com/yaruz/app/internal/pkg/auth"
	"github.com/yaruz/app/internal/pkg/config"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"

	"github.com/yaruz/app/internal/domain/user"

	"github.com/minipkg/log"
	"github.com/minipkg/ozzo_routing/errorshandler"
)

// Service encapsulates the authentication logic.
type Service interface {
	GetSignUpUrl() string
	GetSignInUrl() string
	SignIn(code string, state string) (auth.Claims, error)
	SignUp()
	// authenticate authenticates a user using username and password.
	// It returns a JWT token if authentication succeeds. Otherwise, an error is returned.
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, username, password string) (string, error)
	NewUser(username, password string) (*user.User, error)
	StringTokenValidation(ctx context.Context, stringToken string) (resCtx context.Context, isValid bool, err error)
}

type UserService interface {
}

type service struct {
	//signingKey        string
	//tokenExpiration   uint
	userService user.IService
	logger      log.ILogger
	//sessionRepository Repository
	//tokenRepository   TokenRepository
	Endpoint        string
	ClientId        string
	ClientSecret    string
	Organization    string
	Application     string
	JWTSigningKey   string
	JWTExpiration   uint
	SessionlifeTime uint
}

type contextKey int

const (
	saltSize                  = 64
	iterations                = 1e4
	userSessionKey contextKey = iota
)

// NewService creates a new authentication service.
func NewService(logger log.ILogger, cfg config.Auth, userService user.IService) *service {
	return &service{
		logger:          logger,
		Endpoint:        cfg.Endpoint,
		ClientId:        cfg.ClientId,
		ClientSecret:    cfg.ClientSecret,
		Organization:    cfg.Organization,
		Application:     cfg.Application,
		JWTSigningKey:   cfg.JWTSigningKey,
		JWTExpiration:   cfg.JWTExpiration,
		SessionlifeTime: cfg.SessionlifeTime,
		userService:     userService,
	}
}

func (s service) NewSession(ctx context.Context, userId uint, langId uint) (*session.Session, error) {
	user, err := s.userService.Get(ctx, userId, langId)
	if err != nil {
		return nil, err
	}
	return &session.Session{
		UserID: userId,
		User:   *user,
	}, nil
}

func (s service) NewUser(username, password string) (*user.User, error) {
	user := s.userService.NewEntity()
	user.Name = username

	salt, err := generateRandomBytes(saltSize)
	if err != nil {
		return user, errors.Wrapf(err, "could not get salt: %v", err)
	}
	user.Passhash = string(hashPassword([]byte(password), salt))
	return user, nil
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.authenticate(ctx, username, password)

	if err != nil {
		return "", err
	}

	session, err := s.sessionRepository.Get(ctx, user.ID)
	if err != nil {
		return s.createSession(ctx, *user)
	}
	return s.updateSession(ctx, *user, session)
}

func (s service) updateSession(ctx context.Context, user user.User, sess *session.Session) (string, error) {
	token, err := s.getStringTokenByUser(user)
	if err != nil {
		return "", err
	}
	sess.Token = token
	sess.User = user
	sess.Data = session.Data{
		UserID:              user.ID,
		UserName:            user.Name,
		ExpirationTokenTime: s.getTokenExpirationTime(),
	}

	return token, s.sessionRepository.Update(ctx, sess)
}

func (s service) createSession(ctx context.Context, user user.User) (string, error) {
	token, err := s.getStringTokenByUser(user)
	if err != nil {
		return "", err
	}

	sess, err := s.sessionRepository.NewEntity(ctx, user.ID)
	if err != nil {
		return "", err
	}
	sess.Token = token
	sess.User = user
	sess.Data = session.Data{
		UserID:              user.ID,
		UserName:            user.Name,
		ExpirationTokenTime: s.getTokenExpirationTime(),
	}

	err = s.sessionRepository.Create(ctx, sess)
	if err != nil {
		return "", err
	}

	ctx = context.WithValue(
		ctx,
		userSessionKey,
		sess,
	)
	sess.Ctx = ctx
	return token, nil
}

func (s service) getTokenExpirationTime() time.Time {
	return time.Now().Add(time.Duration(int64(s.tokenExpiration)) * time.Hour)
}

func (s service) getStringTokenByUser(user user.User) (string, error) {
	token := s.tokenRepository.NewTokenByData(auth2.TokenData{
		UserID:              user.ID,
		UserName:            user.Name,
		ExpirationTokenTime: s.getTokenExpirationTime(),
	})
	return token.GenerateStringToken(s.signingKey)
}

// authenticate authenticates a user using username and password.
// If username and password are correct, an *user.User is returned. Otherwise, error is returned.
func (s service) authenticate(ctx context.Context, username, password string) (*user.User, error) {
	logger := s.logger.With(ctx, "user", username)

	user := s.userService.NewEntity()
	user.Name = username

	user, err := s.userService.First(ctx, user)
	if err != nil {
		return user, errorshandler.BadRequest("User not found")
	}

	if comparePassword([]byte(user.Passhash), []byte(password)) {
		logger.Infof("authentication successful")
		return user, nil
	}

	logger.Infof("authentication failed")
	return user, errorshandler.Unauthorized("")
}

func (s service) Register(ctx context.Context, username, password string) (string, error) {
	user, err := s.NewUser(username, password)
	if err != nil {
		return "", errorshandler.InternalServerError(err.Error())
	}

	if err := s.userService.Create(ctx, user); err != nil {
		return "", errorshandler.BadRequest(err.Error())
	}

	return s.createSession(ctx, *user)
}

func (s service) StringTokenValidation(ctx context.Context, stringToken string) (resCtx context.Context, isValid bool, err error) {
	resCtx = ctx
	token, err := s.tokenRepository.ParseStringToken(stringToken, s.signingKey)
	if err != nil {
		return resCtx, isValid, err
	}

	session, err := s.sessionRepository.Get(ctx, token.GetData().UserID)
	if err != nil {
		return resCtx, isValid, err
	}
	isValid = true

	resCtx = context.WithValue(
		ctx,
		userSessionKey,
		session,
	)
	session.Ctx = resCtx
	return resCtx, isValid, nil
}

// Source: https://play.golang.org/p/tAZtO7L6pm
// hash provided clear text password and compare it to provided hash
func comparePassword(hash, pw []byte) bool {
	return bytes.Equal(hash, hashPassword(pw, hash[:saltSize]))
}

// hash the password with the provided salt using the pbkdf2 algorithm
// return byte slice containing salt (first 64 bytes) and hash (last 32 bytes) => total of 96 bytes
func hashPassword(pw, salt []byte) []byte {
	ret := make([]byte, len(salt))
	copy(ret, salt)
	return append(ret, pbkdf2.Key(pw, salt, iterations, sha256.Size, sha256.New)...)
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}
