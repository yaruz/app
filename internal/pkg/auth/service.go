package auth

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/casdoor/casdoor-go-sdk/auth"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/minipkg/log"
	"github.com/pkg/errors"
	"github.com/yaruz/app/internal/domain/account"
	"github.com/yaruz/app/internal/domain/session"
	"github.com/yaruz/app/internal/domain/user"
	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/internal/pkg/config"
	"github.com/yaruz/app/pkg/yarus_platform/reference/domain/text_lang"
	"github.com/yaruz/app/pkg/yarus_platform/yaruserror"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Service encapsulates the authentication logic.
type Service interface {
	GetSignUpUrl() string
	GetSignInUrl() string
	GetForgetUrl() string
	GetSession(ctx context.Context) *session.Session
	SessionInit(ctx context.Context, token string, accountSettings *account.AccountSettings, oauthToken *oauth2.Token) (context.Context, error)
	UpdateSession(ctx context.Context, sess *session.Session) (context.Context, *session.Session, error)
	AccountSettingsUpdate(ctx context.Context, accountSettings *account.AccountSettings) (context.Context, error)
	SignIn(ctx context.Context, code, state string, accountSettings *account.AccountSettings) (context.Context, error)
	StringTokenValidation(ctx context.Context, stringToken string) error
	RoutingGetAccountSettingsWithDefaults(rctx *routing.Context) (*account.AccountSettings, error)
	CheckAuthMiddleware(rctx *routing.Context) error
}

const (
	timeOutInSec          = 5
	URI_RefreshToken      = "/api/login/oauth/refresh_token"
	ctxUserSessionKey     = "session"
	grantTypeRefreshToken = "refresh_token"
	scopeRead             = "read"
)

//go:embed token_jwt_key.pem
var JwtPublicKey string

var _ Service = &service{}

type UserService interface {
}

type service struct {
	//signingKey        string
	//tokenExpiration   uint
	userService user.IService
	logger      log.ILogger
	session     session.Repository
	langFinder  text_lang.LangFinder
	//tokenRepository   TokenRepository
	Endpoint               string
	ClientId               string
	ClientSecret           string
	Organization           string
	Application            string
	SignInRedirectURL      string
	JWTSigningKey          string
	JWTExpiration          uint
	SessionlifeTime        uint
	defaultAccountSettings *account.AccountSettings
}

// NewService creates a new authentication service.
func NewService(ctx context.Context, logger log.ILogger, cfg config.Auth, userService user.IService, session session.Repository, langFinder text_lang.LangFinder) (*service, error) {
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
		langFinder:        langFinder,
	}
	defaultLangID, err := s.langFinder.GetIDByCode(ctx, cfg.DefaultAccountSettings.Lang)
	if err != nil {
		return nil, err
	}

	s.defaultAccountSettings = &account.AccountSettings{
		LangID: defaultLangID,
	}
	// casdoor-go-sdk/auth
	auth.InitConfig(s.Endpoint, s.ClientId, s.ClientSecret, JwtPublicKey, s.Organization, s.Application)
	return s, nil
}

func (s *service) RoutingGetAccountSettingsWithDefaults(rctx *routing.Context) (*account.AccountSettings, error) {
	ctx := rctx.Request.Context()
	accountSettings := s.defaultAccountSettings

	if lang := rctx.Request.Header.Get("langID"); lang != "" {
		langID, err := s.langFinder.GetIDByCode(ctx, lang)
		if err != nil {
			if err != yaruserror.ErrNotFound {
				return nil, err
			}
			// todo: временно игнорим ненайденные
		} else {
			accountSettings.LangID = langID
		}
	}

	return accountSettings, nil
}

func (s *service) newSession(ctx context.Context, jwtClaims *auth.Claims, user *user.User, accountSettings *account.AccountSettings, oauthToken *oauth2.Token) *session.Session {
	return &session.Session{
		User:            user,
		JwtClaims:       jwtClaims,
		Token:           oauthToken,
		AccountSettings: accountSettings,
	}
}

func (s *service) createSession(ctx context.Context, jwtClaims *auth.Claims, user *user.User, defaultAccountSettings *account.AccountSettings, oauthToken *oauth2.Token) (context.Context, *session.Session, error) {
	if jwtClaims == nil || user == nil || defaultAccountSettings == nil || oauthToken == nil {
		return ctx, nil, errors.Wrapf(apperror.ErrBadParams, "jwtClaims == %v \nuser == %v \ndefaultAccountSettings == %v \noauthToken == %v", jwtClaims, user, defaultAccountSettings, oauthToken)
	}

	if accSettings, err := s.accountGetSettings(&jwtClaims.User); err != nil && accSettings != nil {
		defaultAccountSettings = accSettings
	}

	sess := s.newSession(ctx, jwtClaims, user, defaultAccountSettings, oauthToken)

	err := s.session.Create(ctx, sess)
	if err != nil {
		return ctx, nil, err
	}

	ctx = context.WithValue(
		ctx,
		ctxUserSessionKey,
		sess,
	)

	return ctx, sess, nil
}

func (s *service) UpdateSession(ctx context.Context, sess *session.Session) (context.Context, *session.Session, error) {
	err := s.session.Update(ctx, sess)
	if err != nil {
		return ctx, nil, err
	}

	ctx = context.WithValue(
		ctx,
		ctxUserSessionKey,
		sess,
	)

	return ctx, sess, nil
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
	return ctx.Value(ctxUserSessionKey).(*session.Session)
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

func (s *service) accountGetSettings(acc *auth.User) (*account.AccountSettings, error) {
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

func (s *service) SessionInit(ctx context.Context, accessToken string, accountSettings *account.AccountSettings, oauthToken *oauth2.Token) (context.Context, error) {
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

	//affected, err := object.UpdateMemberOnlineStatus(&claims.User, true, util.GetCurrentTime())
	//if err != nil {
	//	c.ResponseError(err.Error())
	//	return
	//}

	return ctx, nil
}

func (s *service) accountPropertiesUpdate(ctx context.Context, refreshToken string, account *auth.User) (*oauth2.Token, *auth.Claims, error) {
	ok, err := auth.UpdateUserForColumns(account, []string{"properties"})
	if !ok || err != nil {
		return nil, nil, errors.Wrapf(apperror.ErrInternal, fmt.Sprintf("User not updated. Ok = %t, err = %q.", ok, err.Error()))
	}

	return s.refreshAndParseToken(ctx, refreshToken)

}

func (s *service) getAndParseToken(ctx context.Context, code, state string) (*oauth2.Token, *auth.Claims, error) {
	authToken, err := auth.GetOAuthToken(code, state)
	if err != nil {
		return nil, nil, err
	}

	jwtClaims, err := auth.ParseJwtToken(authToken.AccessToken)
	if err != nil {
		return nil, nil, err
	}
	jwtClaims.AccessToken = authToken.AccessToken

	return authToken, jwtClaims, nil
}

type refreshTokenParams struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type refreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	Error        string `json:"error"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func (s *service) refreshToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	httpClient := &http.Client{
		//Timeout: time.Second * 10,
	}

	urlValues := url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{refreshToken},
		"scope":         []string{"read"},
		"client_id":     []string{s.ClientId},
		"client_secret": []string{s.ClientSecret},
	}

	body := bytes.NewBuffer([]byte{})

	reqUrl := s.Endpoint + URI_RefreshToken + "?" + urlValues.Encode()
	req, err := http.NewRequest(http.MethodPost, reqUrl, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(body.Len()))

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res refreshTokenResponse
	if err = json.Unmarshal(respBody, &res); err != nil {
		return nil, err
	}

	expires := res.ExpiresIn

	newOauthToken := &oauth2.Token{
		AccessToken:  res.AccessToken,
		TokenType:    res.TokenType,
		RefreshToken: res.RefreshToken,
		Expiry:       time.Now().Add(time.Duration(expires) * time.Second), // должно быть в секундах, но в Касдоре ошибка
	}

	if strings.HasPrefix(newOauthToken.AccessToken, "error:") {
		return nil, errors.New(strings.TrimLeft(newOauthToken.AccessToken, "error: "))
	}

	return newOauthToken, err
}

func (s *service) autoRefreshToken(ctx context.Context, oauthToken *oauth2.Token) (*oauth2.Token, error) {
	config := &oauth2.Config{
		ClientID:     s.ClientId,
		ClientSecret: s.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   fmt.Sprintf("%s/api/login/oauth/authorize", s.Endpoint),
			TokenURL:  fmt.Sprintf("%s/api/login/oauth/access_token", s.Endpoint),
			AuthStyle: oauth2.AuthStyleInParams,
		},
		Scopes: nil,
	}
	tokenSource := config.TokenSource(ctx, oauthToken)

	newOauthToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(newOauthToken.AccessToken, "error:") {
		return nil, errors.New(strings.TrimLeft(newOauthToken.AccessToken, "error: "))
	}

	return newOauthToken, err
}

func (s *service) refreshAndParseToken(ctx context.Context, refreshToken string) (*oauth2.Token, *auth.Claims, error) {
	//oauthToken, err := RefreshToken(ctx, config, refreshToken)
	oauthToken, err := s.refreshToken(ctx, refreshToken)
	if err != nil {
		return nil, nil, err
	}

	jwtClaims, err := auth.ParseJwtToken(oauthToken.AccessToken)
	if err != nil {
		return nil, nil, err
	}
	jwtClaims.AccessToken = oauthToken.AccessToken

	return oauthToken, jwtClaims, nil
}

func (s *service) signUp(ctx context.Context, jwtClaims *auth.Claims, langId uint) (*user.User, error) {
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

func (s *service) StringTokenValidation(ctx context.Context, stringToken string) error {
	//	temporary validation method
	_, err := auth.ParseJwtToken(stringToken)
	if err != nil {
		return err
	}
	return nil
}

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
