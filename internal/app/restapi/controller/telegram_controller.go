package controller

import (
	"context"
	"fmt"

	"github.com/gotd/td/tg"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/gotd/td/session"
	"github.com/yaruz/app/internal/pkg/apperror"
	pkg_tg "github.com/yaruz/app/internal/pkg/socnets/tg"
	"sync"

	"github.com/minipkg/log"

	"github.com/yaruz/app/internal/domain/user"
	"github.com/yaruz/app/internal/pkg/auth"
)

const (
	appID   = int(16433962)
	appHash = "c47038063d593134511bb756011fb409"
)

type TgService interface {
	IsAuth(ctx context.Context, sessionID string) (bool, error)
	SendCode(ctx context.Context, sessionID string, phone string) error
	SignIn(ctx context.Context, sessionID string, code string, password string) (*tg.User, error)
}

// memorySession implements in-memory session storage.
// Goroutine-safe.
type memorySession struct {
	mux  sync.RWMutex
	data []byte
}

// LoadSession loads session from memory.
func (s *memorySession) LoadSession(context.Context) ([]byte, error) {
	if s == nil {
		return nil, session.ErrNotFound
	}

	s.mux.RLock()
	defer s.mux.RUnlock()

	if len(s.data) == 0 {
		return nil, session.ErrNotFound
	}

	cpy := append([]byte(nil), s.data...)

	return cpy, nil
}

// StoreSession stores session to memory.
func (s *memorySession) StoreSession(ctx context.Context, data []byte) error {
	s.mux.Lock()
	s.data = data
	s.mux.Unlock()
	return nil
}

type telegramController struct {
	RouteGroup  *routing.RouteGroup
	logger      log.Logger
	userService user.IService
	authService auth.Service
	tgService   TgService
}

func NewTelegramController(r *routing.RouteGroup, logger log.Logger, authService auth.Service, userService user.IService, tgService pkg_tg.Service) *telegramController {
	return &telegramController{
		RouteGroup:  r,
		logger:      logger,
		userService: userService,
		authService: authService,
		tgService:   tgService,
	}
}

// RegisterHandlers sets up the routing of the HTTP handlers.
func (c *telegramController) RegisterHandlers() {

	c.RouteGroup.Get(`/is-auth`, c.isAuth)
	c.RouteGroup.Get(`/send-code/<phone>`, c.sendCode)
	//c.RouteGroup.Get(`/sign-in/<code>`, c.signIn)
	c.RouteGroup.Post(`/sign-in/<code>`, c.signIn)
	//c.RouteGroup.Post(`/sess-pass`, c.authCheckSessionPassword)

}

func (c *telegramController) isAuth(rctx *routing.Context) (err error) {
	ctx := rctx.Request.Context()

	sess, err := c.authService.GetSession(ctx)
	if err != nil {
		return err
	}

	isAuth, err := c.tgService.IsAuth(ctx, sess.ID)
	if err != nil {
		return err
	}
	return rctx.Write(isAuth)
}

func (c *telegramController) sendCode(rctx *routing.Context) (err error) {
	ctx := rctx.Request.Context()

	phone := rctx.Param("phone")
	if phone == "" {
		return fmt.Errorf("[%w] Phone is empty.", apperror.ErrBadParams)
	}

	sess, err := c.authService.GetSession(ctx)
	if err != nil {
		return err
	}

	if err = c.tgService.SendCode(ctx, sess.ID, phone); err != nil {
		return err
	}
	return rctx.Write(true)
}

func (c *telegramController) signIn(rctx *routing.Context) (err error) {
	ctx := rctx.Request.Context()

	code := rctx.Param("code")
	if code == "" {
		return fmt.Errorf("[%w] Code is empty.", apperror.ErrBadParams)
	}
	var pass string
	err = rctx.Read(&pass)

	sess, err := c.authService.GetSession(ctx)
	if err != nil {
		return err
	}

	tgUser, err := c.tgService.SignIn(ctx, sess.ID, code, pass)
	if err != nil {
		return err
	}
	return rctx.Write(*tgUser)
}

/*
// authSendCode сохраняет телефон пользователя и отправляет ему в Телеграм код аутентификации
func (c *telegramController) authSendCode(rctx *routing.Context) (err error) {
	ctx := rctx.Request.Context()
	sess := c.authService.GetSession(ctx)
	//isTgAuthSessionRegistred, err := c.tgService.IsAuthSessionRegistred(sess)
	//if err != nil {
	//	return err
	//}
	//if isTgAuthSessionRegistred {
	//	return rctx.Write("You've already signed in!")
	//}

	phone := rctx.Param("phone")
	if phone == "" {
		return errors.Wrapf(apperror.ErrBadParams, "Phone must be set.")
	}
	// Сохраняем телефон в User
	sess.User, err = c.User.Get(ctx, sess.User.ID, sess.AccountSettings.LangID)
	if err != nil {
		return err
	}
	if err = sess.User.SetPhone(ctx, phone); err != nil {
		return err
	}

	if err = c.User.Update(ctx, sess.User, sess.AccountSettings.LangID); err != nil {
		return nil
	}
	// Сохраняем телефон в аккаунте, обновляем и сохраняем сессию
	sess.JwtClaims.User.Phone = phone
	ctx, err = c.authService.AccountUpdate(ctx, sess)
	if err != nil {
		return err
	}
	*rctx.Request = *rctx.Request.WithContext(ctx)

	return rctx.Write(true)
}

// authSignIn - аутентификация в Телеграм по отправленному коду Телеграм
func (c *telegramController) authSignIn(rctx *routing.Context) error {
	ctx := rctx.Request.Context()
	sess := c.authService.GetSession(ctx)
	isAuthSessionRegistred, err := c.tgService.IsAuthSessionRegistred(sess)
	if err != nil {
		return err
	}
	if isAuthSessionRegistred {
		return rctx.Write("You've already signed in!")
	}

	code := rctx.Param("code")
	if code == "" {
		return errors.Wrapf(apperror.ErrBadParams, "Code must be set.")
	}

	client, err := c.tgService.NewClient(sess)
	if err != nil {
		return err
	}
	// Запрос аутентификации в Телеграм
	if err := c.tgService.AuthSignIn(ctx, client, sess.User.Phone, code); err != nil {
		return err
	}
	return rctx.Write(true)
}

// authCheckSessionPassword - проверка сессионного пороля от аккаунта Телеграм
// Данный пароль нигде не сохраняется и никуда не пересылается, а используется только в математическом алгоритме проверки подлинности владельца аккаунта Телеграм.
func (c *telegramController) authCheckSessionPassword(rctx *routing.Context) error {
	ctx := rctx.Request.Context()
	sess := c.authService.GetSession(ctx)
	isAuthSessionRegistred, err := c.tgService.IsAuthSessionRegistred(sess)
	if err != nil {
		return err
	}
	if isAuthSessionRegistred {
		return rctx.Write("You've already signed in!")
	}

	client, err := c.tgService.NewClient(sess)
	if err != nil {
		return err
	}

	var pass string
	rctx.Read(&pass)
	if pass == "" {
		return errors.Wrapf(apperror.ErrBadParams, "Session password for Telegram must be set.")
	}

	if err := c.tgService.AuthCheckPassword(ctx, client, pass); err != nil {
		return err
	}
	return rctx.Write(true)
}

*/
