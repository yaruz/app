package auth

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/minipkg/log"
	"github.com/yaruz/app/internal/domain/tg_account"
	"github.com/yaruz/app/internal/pkg/jwt"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_lang"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
	"net/http"
	"time"

	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/internal/pkg/session"

	"github.com/yaruz/app/internal/domain/account"
	"github.com/yaruz/app/internal/domain/user"
)

// Service encapsulates the authentication logic.
type Service interface {
	//GetSignUpUrl() string
	//GetSignInUrl() string
	//GetForgetUrl() string
	//GetSession(ctx context.Context) *session.Session
	//SessionInit(ctx context.Context, token string, sessionID string, accountSettings *account.AccountSettings) (context.Context, error)
	//UpdateSession(ctx context.Context, sess *session.Session) (context.Context, *session.Session, error)
	//AccountSettingsUpdate(ctx context.Context, accountSettings *account.AccountSettings) (context.Context, error)
	//AccountUpdate(ctx context.Context, sess *session.Session) (context.Context, error)
	//SignIn(ctx context.Context, code, state string, accountSettings *account.AccountSettings) (context.Context, error)
	//StringTokenValidation(ctx context.Context, stringToken string) error
	//RoutingGetAccountSettingsWithDefaults(rctx *routing.Context) (*account.AccountSettings, error)
	//CheckAuthMiddleware(rctx *routing.Context) error
	//CreateUnauthSession(ctx context.Context, accountSettings *account.AccountSettings) (*session.Session, error)
	GetSession(ctx context.Context) (*session.Session, error)
	MiddlewareCheckAuth() routing.Handler
	MiddlewareSessionInit() routing.Handler
	SignInTg(ctx context.Context, account *tg_account.TgAccount, langID uint) (accessToken string, err error)
}

const (
	timeOutInSec     = 5
	URI_RefreshToken = "/api/login/oauth/refresh_token"
	ctx_SessionIDKey = "session"
	ctx_IsAuth       = "isAuth"
)

//go_:embed token_jwt_key.pem
//var JwtPublicKey string

var _ Service = (*service)(nil)

type UserService interface {
}

type service struct {
	userService            user.IService
	logger                 log.Logger
	sessionRepository      session.Repository
	langFinder             text_lang.LangFinder
	tokenRepository        TokenRepository
	SignInRedirectURL      string
	SessionlifeTimeInHours uint
	defaultAccountSettings *account.AccountSettings
}

type Config struct {
	SignInRedirectURL      string
	JWTSigningKey          string
	JWTExpirationInHours   uint
	SessionlifeTimeInHours uint
	DefaultAccountSettings DefaultAccountSettings
}

type DefaultAccountSettings struct {
	Lang string
}

// NewService creates a new authentication service.
func NewService(ctx context.Context, logger log.Logger, cfg Config, tokenRepository TokenRepository, sessionRepository session.Repository, userService user.IService, langFinder text_lang.LangFinder) (*service, error) {
	s := &service{
		logger:                 logger,
		SignInRedirectURL:      cfg.SignInRedirectURL,
		SessionlifeTimeInHours: cfg.SessionlifeTimeInHours,
		tokenRepository:        tokenRepository,
		userService:            userService,
		sessionRepository:      sessionRepository,
		langFinder:             langFinder,
	}
	defaultLangID, err := s.langFinder.GetIDByCode(ctx, cfg.DefaultAccountSettings.Lang)
	if err != nil {
		return nil, err
	}

	s.defaultAccountSettings = &account.AccountSettings{
		LangID: defaultLangID,
	}
	return s, nil
}

func (s *service) GetSession(ctx context.Context) (*session.Session, error) {
	sessionID, err := s.getSessionID(ctx)
	if err != nil {
		return nil, fmt.Errorf("[%w] sessionID not found in context", err)
	}
	return s.sessionRepository.Get(ctx, sessionID)
}

func (s *service) SaveSession(ctx context.Context, entity *session.Session) error {
	return s.sessionRepository.Set(ctx, entity)
}

func (s *service) routingGetAccountSettingsWithDefaults(rctx *routing.Context) (*account.AccountSettings, error) {
	ctx := rctx.Request.Context()
	accountSettings := s.defaultAccountSettings

	if lang := rctx.Request.Header.Get("langID"); lang != "" {
		langID, err := s.langFinder.GetIDByCode(ctx, lang)
		if err != nil {
			if err != yaruserror.ErrNotFound {
				return nil, fmt.Errorf("[%w] langID error: %s", apperror.ErrInternal, err.Error())
			}
			return nil, fmt.Errorf("[%w] langID error: %s", apperror.ErrBadParams, err.Error())
		} else {
			accountSettings.LangID = langID
		}
	}

	return accountSettings, nil
}

func (s *service) stringTokenValidation(ctx context.Context, stringToken string) error {
	//	temporary validation method
	_, err := s.tokenRepository.ParseStringToken(stringToken)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) MiddlewareCheckAuth() routing.Handler {
	return func(rctx *routing.Context) error {
		ctx := rctx.Request.Context()
		val := ctx.Value(ctx_IsAuth)
		isAuth, ok := val.(bool)
		if !ok {
			return routing.NewHTTPError(http.StatusInternalServerError, "Context has been broken.")
		}

		if !isAuth {
			return routing.NewHTTPError(http.StatusUnauthorized, "")
		}
		return nil
	}
}

func (s *service) MiddlewareSessionInit() routing.Handler {
	return func(rctx *routing.Context) error {
		// todo: все сообщения при ошибках нужно выдавать не наружу, а в логи
		ctx := rctx.Request.Context()
		var isAuth bool
		var err error
		var sessionID string

		sessionID = getSessionID(rctx.Request.Header.Get(header_Cookie))
		stringToken := getAuthorizationToken(rctx.Request.Header.Get(header_Authorization))

		accountSettings, err := s.routingGetAccountSettingsWithDefaults(rctx)
		if err != nil {
			return routing.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if stringToken != "" {
			// если есть токен, проверяем токен, достаём сессию, аутентифицируем пользователя
			sess, err := s.authSessionInit(ctx, stringToken)
			if err != nil {
				if !errors.Is(err, apperror.ErrTokenHasExpired) {
					return routing.NewHTTPError(http.StatusUnauthorized, err.Error())
				}
				return routing.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			// обновим в сессии accountSettings
			sess.AccountSettings = accountSettings
			if err = s.sessionRepository.Set(ctx, sess); err != nil {
				return routing.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			sessionID = sess.ID
			isAuth = true
		} else if sessionID != "" {
			// если, нет токена, но есть sessID - достаём сессию для неаутентифицированного пользователя
			sess, ok := s.unauthCheckSessionID(ctx, sessionID)
			if ok {
				// обновим в сессии accountSettings
				sess.AccountSettings = accountSettings
				if err = s.sessionRepository.Set(ctx, sess); err != nil {
					return routing.NewHTTPError(http.StatusInternalServerError, err.Error())
				}

				sessionID = sess.ID
			} else {
				sessionID = ""
			}
		}

		if sessionID == "" {
			// нет ни какой сессии, создаём новую сессию для неаутентифицированного пользователя
			sess, err := s.unauthSessionCreate(ctx, accountSettings)
			if err != nil {
				return routing.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			sessionID = sess.ID
		}

		ctx = context.WithValue(
			ctx,
			ctx_SessionIDKey,
			sessionID,
		)
		ctx = context.WithValue(
			ctx,
			ctx_IsAuth,
			isAuth,
		)
		*rctx.Request = *rctx.Request.WithContext(ctx)
		rctx.Response.Header().Set(header_Cookie, cookie_SessionID+"="+sessionID)
		return nil
	}
}

func (s *service) authSessionInit(ctx context.Context, accessToken string) (*session.Session, error) {
	token, err := s.tokenRepository.ParseStringToken(accessToken)
	if err != nil {
		return nil, err
	}

	tokenData := token.GetData()
	sess, ok := s.checkAuthData(ctx, tokenData.SessionID, tokenData.User.ID, accessToken)
	if !ok {
		return nil, apperror.ErrInvalidToken
	}

	return sess, nil
}

func (s *service) checkAuthData(ctx context.Context, sessionID string, userID uint, accessToken string) (*session.Session, bool) {
	sess, err := s.sessionRepository.Get(ctx, sessionID)
	return sess, err == nil && sess.JwtClaims.AccessToken == accessToken && sess.User.ID == userID
}

func (s *service) unauthCheckSessionID(ctx context.Context, sessionID string) (*session.Session, bool) {
	sess, err := s.sessionRepository.Get(ctx, sessionID)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil, false
		}
		s.logger.Errorf("Internal server error: %s", err.Error())
		return nil, false
	}
	return sess, true
}

func (s *service) unauthSessionCreate(ctx context.Context, accountSettings *account.AccountSettings) (*session.Session, error) {
	sess, err := session.New("", accountSettings, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("[%w] Session creation error.", apperror.ErrInternal)
	}
	err = s.sessionRepository.Set(ctx, sess)
	if err != nil {
		return nil, fmt.Errorf("[%w] Session save error.", apperror.ErrInternal)
	}
	return sess, nil
}

func (s *service) getSessionID(ctx context.Context) (string, error) {
	sessionID := ctx.Value(ctx_SessionIDKey).(string)
	if sessionID == "" {
		return "", apperror.ErrNotFound
	}
	return sessionID, nil
}

func (s *service) SignInTg(ctx context.Context, account *tg_account.TgAccount, langID uint) (accessToken string, err error) {
	// проверить, мож, уже есть
	u, err := s.userService.GetByTgAccount(ctx, account, langID)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return "", fmt.Errorf("[%w] get user by tg account error: %s", apperror.ErrInternal, err.Error())
		}

		u, err = s.createUserTg(ctx, account)
		if err != nil {
			return "", err
		}
	}
	return s.signIn(ctx, u)
}

func (s *service) signIn(ctx context.Context, u *user.User) (accessToken string, err error) {
	sess, err := s.GetSession(ctx)
	if err != nil {
		return "", err
	}

	token, err := s.tokenRepository.NewTokenWithData(&jwt.TokenData{
		SessionID: sess.ID,
		User:      u,
	})
	if err != nil {
		return "", err
	}

	sess.User = u
	sess.JwtClaims = token.GetClaims()

	accessToken, err = token.GenerateStringToken()
	if err != nil {
		return "", fmt.Errorf("[%w] generate string token error: %s", apperror.ErrInternal, err.Error())
	}

	return accessToken, s.sessionRepository.Set(ctx, sess)
}

func (s *service) createUserTg(ctx context.Context, account *tg_account.TgAccount) (*user.User, error) {
	sess, err := s.GetSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("[%w] getting session error: %s", apperror.ErrInternal, err.Error())
	}

	u, err := s.userService.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("[%w] user creation error: %s", apperror.ErrInternal, err.Error())
	}

	err = u.SetPhone(ctx, account.Phone)
	if err != nil {
		return nil, fmt.Errorf("[%w] setting phone error: %s", apperror.ErrInternal, err.Error())
	}
	err = u.SetFirstName(ctx, account.FirstName)
	if err != nil {
		return nil, fmt.Errorf("[%w] setting first name error: %s", apperror.ErrInternal, err.Error())
	}
	err = u.SetLastName(ctx, account.LastName)
	if err != nil {
		return nil, fmt.Errorf("[%w] setting last name error: %s", apperror.ErrInternal, err.Error())
	}
	err = u.SetUserName(ctx, account.UserName)
	if err != nil {
		return nil, fmt.Errorf("[%w] setting user name error: %s", apperror.ErrInternal, err.Error())
	}
	err = u.SetCreatedAt(ctx, time.Now().UTC())
	if err != nil {
		return nil, fmt.Errorf("[%w] setting time of creation error: %s", apperror.ErrInternal, err.Error())
	}

	if err = s.userService.Create(ctx, u, sess.AccountSettings.LangID); err != nil {
		return nil, fmt.Errorf("[%w] user saving error: %s", apperror.ErrInternal, err.Error())
	}

	return u, nil
}

/*
func (s *service) saveSession(ctx context.Context, sess *session.Session) (context.Context, error) {

	if err := s.sessionRepository.Create(ctx, sess); err != nil {
		return ctx, err
	}

	return context.WithValue(
		ctx,
		ctxUser_SessionKey,
		sess,
	), nil
}

func (s *service) UpdateSession(ctx context.Context, sess *session.Session) (context.Context, *session.Session, error) {
	err := s.sessionRepository.Update(ctx, sess)
	if err != nil {
		return ctx, nil, err
	}

	ctx = context.WithValue(
		ctx,
		ctxUser_SessionKey,
		sess,
	)

	return ctx, sess, nil
}

func (s *service) CreateUnauthSession(ctx context.Context, accountSettings *account.AccountSettings) (context.Context, *session.Session, error) {
	sess, err := session.New("", accountSettings, nil, nil)
	if err != nil {
		return ctx, nil, err
	}

	ctx, err = s.saveSession(ctx, sess)
	return ctx, sess, err
}

func (s *service) createSession(ctx context.Context, jwtClaims *jwt.Claims, user *user.User, defaultAccountSettings *account.AccountSettings) (context.Context, *session.Session, error) {
	if jwtClaims == nil || user == nil || defaultAccountSettings == nil {
		return ctx, nil, errors.Wrapf(apperror.ErrBadParams, "jwtClaims: %v \nuser: %v \ndefaultAccountSettings: %v \n", jwtClaims, user, defaultAccountSettings)
	}

	if accSettings, err := s.accountGetSettings(jwtClaims.User); err != nil && accSettings != nil {
		defaultAccountSettings = accSettings
	}

	sess, err := session.New("", defaultAccountSettings, jwtClaims, user)
	if err != nil {
		return ctx, nil, err
	}

	ctx, err = s.saveSession(ctx, sess)
	return ctx, sess, err
}

func (s *service) updateSession(ctx context.Context, sess *session.Session, jwtClaims *auth.Claims, user *user.User, accountSettings *account.AccountSettings, oauthToken *oauth2.Token) (context.Context, *session.Session, error) {
	if user != nil {
		sess.User = user
	}

	if accountSettings != nil {
		sess.AccountSettings = accountSettings
	}

	if jwtClaims != nil {
		sess.JwtClaims = jwtClaims
	}

	if oauthToken != nil {
		sess.Token = oauthToken
	}

	return s.UpdateSession(ctx, sess)
}

func (s *service) GetSession(ctx context.Context) *session.Session {
	return ctx.Value(ctxUser_SessionKey).(*session.Session)
}

func (s *service) GetSignUpUrl() string {
	return fmt.Sprintf("%s/signup/%s", s.Endpoint, s.Application)
}

func (s *service) GetSignInUrl() string {
	return fmt.Sprintf("%s/login/oauth/authorize?client_id=%s&response_type=code&scope=read&state=%s&redirect_uri=%s", s.Endpoint, s.ClientId, s.Application, s.SignInRedirectURL)
}

func (s *service) GetForgetUrl() string {
	return fmt.Sprintf("%s/forget/%s", s.Endpoint, s.Application)
}

func (s *service) accountGetUintParam(account *auth.User, paramName string) (uint, error) {
	if len(account.Properties) == 0 {
		return 0, apperror.ErrNotFound
	}
	pStr, ok := account.Properties[paramName]
	if !ok {
		return 0, apperror.ErrNotFound
	}

	pInt, err := strconv.Atoi(pStr)
	if err != nil {
		return 0, err
	}

	return uint(pInt), nil
}

func (s *service) accountSetUintParam(account *auth.User, key string, value uint) {
	if account.Properties == nil {
		account.Properties = make(map[string]string, 1)
	}
	account.Properties[key] = strconv.Itoa(int(value))
}

func (s *service) accountSetUserID(account *auth.User, userID uint) {
	s.accountSetUintParam(account, s.Application+"Id", userID)
}

func (s *service) accountGetUserID(account *auth.User) (uint, error) {
	userID, err := s.accountGetUintParam(account, s.Application+"Id")
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *service) accountSetSettings(acc *auth.User, accountSettings *account.AccountSettings) {
	acc.Properties["langId"] = strconv.Itoa(int(accountSettings.LangID))
}

func (s *service) accountGetSettings(acc *user.User) (*account.AccountSettings, error) {
	langID, err := s.accountGetUintParam(acc, "langId")
	if err != nil {
		return nil, err
	}

	return &account.AccountSettings{
		LangID: langID,
	}, nil
}

func (s *service) AccountSettingsUpdate(ctx context.Context, accountSettings *account.AccountSettings) (context.Context, error) {
	sess := s.GetSession(ctx)
	oldAccountSettings, err := s.accountGetSettings(&sess.JwtClaims.User)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return ctx, err
		}
	} else if oldAccountSettings.LangID == accountSettings.LangID { // Если нет изменений - выходим
		return ctx, nil
	}

	s.accountSetSettings(&sess.JwtClaims.User, accountSettings)

	oauthToken, jwtClaims, err := s.accountPropertiesUpdate(ctx, sess.Token.RefreshToken, &sess.JwtClaims.User)
	if err != nil {
		return ctx, err
	}

	ctx, _, err = s.updateSession(ctx, sess, jwtClaims, nil, accountSettings, oauthToken)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (s *service) AccountUpdate(ctx context.Context, sess *session.Session) (context.Context, error) {
	if _, err := auth.UpdateUser(&sess.JwtClaims.User); err != nil {
		return ctx, err
	}
	oauthToken, jwtClaims, err := s.refreshAndParseToken(ctx, sess.Token.RefreshToken)
	if err != nil {
		return ctx, err
	}

	ctx, _, err = s.updateSession(ctx, sess, jwtClaims, nil, nil, oauthToken)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (s *service) SessionInit(ctx context.Context, accessToken string, sessionID string, accountSettings *account.AccountSettings) (context.Context, error) {
	jwtClaims, err := auth.ParseJwtToken(accessToken)
	if err != nil {
		return ctx, err
	}
	jwtClaims.AccessToken = accessToken

	userID, err := s.accountGetUserID(&jwtClaims.User)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return ctx, err
		}
		return ctx, errors.Wrapf(apperror.ErrNotFound, "Param %q not found in JwtClaims.User = %v", s.Application+"ID", jwtClaims.User)
	}

	sess, err := s.sessionRepository.Get(ctx, userID)
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

		tgAcc, err := s.userService.GetTgAccount(ctx, user, accountSettings.LangID)
		if err != nil {
			if !errors.Is(err, apperror.ErrNotFound) {
				return ctx, err
			}
		} else {
			sess.TgAccount = tgAcc
		}

		ctx, _, err = s.createSession(ctx, jwtClaims, user, accountSettings, oauthToken)
		if err != nil {
			return ctx, err
		}
	} else {
		// todo: обновлять только при изменении
		ctx, _, err = s.updateSession(ctx, sess, jwtClaims, nil, nil, oauthToken)
		if err != nil {
			return ctx, err
		}
	}

	return ctx, nil
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s *service) SignIn(ctx context.Context, code, state string, defaultAccountSettings *account.AccountSettings) (context.Context, error) {
	var accountSettings *account.AccountSettings = defaultAccountSettings
	var thereIsNoSavedAccountSettings = true

	oauthToken, jwtClaims, err := s.getAndParseToken(ctx, code, state)
	if err != nil {
		return ctx, err
	}
	// если есть у аккаунта, берём настройки аккаунта
	if accSettings, err := s.accountGetSettings(&jwtClaims.User); err == nil {
		accountSettings = accSettings
		thereIsNoSavedAccountSettings = false
	}
	// получаем user, если нет - создаём нового
	user, err := s.userService.GetByAccountID(ctx, jwtClaims.User.Id, accountSettings.LangID)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return ctx, err
		}
		if user, err = s.signUp(ctx, jwtClaims, accountSettings.LangID); err != nil {
			return nil, err
		}
	}
	// привязываем аккаунт, если не привязан
	userID, err := s.accountGetUserID(&jwtClaims.User)
	if err != nil {
		if !errors.Is(err, apperror.ErrNotFound) {
			return ctx, err
		}

		s.accountSetUserID(&jwtClaims.User, user.ID)

		if oauthToken, jwtClaims, err = s.accountPropertiesUpdate(ctx, oauthToken.RefreshToken, &jwtClaims.User); err != nil {
			return ctx, err
		}
		//	user был создан, но был удалён без удаления аккаунта. Пока не понятно, что с этим делать..
	} else if userID != user.ID {
		return ctx, err
	}

	ctx, err = s.SessionInit(ctx, jwtClaims.AccessToken, accountSettings, oauthToken)

	if thereIsNoSavedAccountSettings {
		sess := s.GetSession(ctx)
		s.AccountSettingsUpdate(ctx, sess.AccountSettings)
	}

	return ctx, nil
}
*/
/*
func (s *service) CheckAuthMiddleware(rctx *routing.Context) error {
	// todo: все сообщения при ошибках нужно выдавать не наружу, а в логи
	ctx := rctx.Request.Context()
	header := rctx.Request.Header.Get("Authorization")
	token := ""
	var err error

	if strings.HasPrefix(header, "Bearer ") {
		token = header[7:]
	} else {
		return UnauthorizedError(rctx, "")
	}

	if err = s.StringTokenValidation(ctx, token); err != nil {
		return UnauthorizedError(rctx, err.Error())
	}

	accountSettings, err := s.RoutingGetAccountSettingsWithDefaults(rctx)
	if err != nil {
		return routing.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if ctx, err = s.SessionInit(ctx, token, accountSettings, nil); err != nil {
		return UnauthorizedError(rctx, err.Error())
	}

	*rctx.Request = *rctx.Request.WithContext(ctx)
	return nil
}
*/
