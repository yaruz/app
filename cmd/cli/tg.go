package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"github.com/minipkg/log"
	"golang.org/x/crypto/ssh/terminal"
	golog "log"
	"os"
	"strings"
	"sync"

	"github.com/yaruz/app/internal/pkg/config"
)

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

var sessionStorage memorySession

// noSignUp can be embedded to prevent signing up.
type noSignUp struct{}

func (c noSignUp) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, errors.New("not implemented")
}

func (c noSignUp) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

// termAuth implements authentication via terminal.
type termAuth struct {
	noSignUp

	phone string
}

func (a termAuth) Phone(_ context.Context) (string, error) {
	return a.phone, nil
}

func (a termAuth) Password(_ context.Context) (string, error) {
	fmt.Print("Enter 2FA password: ")
	bytePwd, err := terminal.ReadPassword(0)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePwd)), nil
}

func (a termAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter code: ")
	//code := "17149"
	code, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(code), nil
}

func main() {

	appID := int(16433962)
	appHash := "c47038063d593134511bb756011fb409"
	phone := "79778508041"

	ctx := context.Background()
	cfg, err := config.Get()
	if err != nil {
		golog.Fatalln("Can not load the config")
	}
	logger, err := log.New(cfg.Infrastructure.Log)
	if err != nil {
		golog.Fatal(err)
	}

	//publicKeys := filepath.Join(appStorage, "tg_public_keys.pem")

	// Setting up authentication flow helper based on terminal auth.
	flow := auth.NewFlow(
		termAuth{phone: phone},
		auth.SendCodeOptions{},
	)

	client := telegram.NewClient(
		appID,
		appHash,
		telegram.Options{
			SessionStorage: &sessionStorage,
			Logger:         logger.ZapLogger(),
		},
	)

	if err != nil {
		logger.Fatal(err)
	}

	err = client.Run(ctx, func(ctx context.Context) error {
		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return err
		}

		logger.Info("Success")

		return nil
	})

	main2(ctx, cfg)

}

func main2(ctx context.Context, cfg *config.Configuration) {
	ctx = context.Background()
	appID := int(16433962)
	appHash := "c47038063d593134511bb756011fb409"
	phone := "79778508041"

	logger, err := log.New(cfg.Infrastructure.Log)
	if err != nil {
		golog.Fatal(err)
	}

	//publicKeys := filepath.Join(appStorage, "tg_public_keys.pem")

	// Setting up authentication flow helper based on terminal auth.
	flow := auth.NewFlow(
		termAuth{phone: phone},
		auth.SendCodeOptions{},
	)

	client := telegram.NewClient(
		appID,
		appHash,
		telegram.Options{
			SessionStorage: &sessionStorage,
			Logger:         logger.ZapLogger(),
		},
	)

	if err != nil {
		logger.Fatal(err)
	}

	err = client.Run(ctx, func(ctx context.Context) error {
		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return err
		}

		logger.Info("Success")

		return nil
	})

}
