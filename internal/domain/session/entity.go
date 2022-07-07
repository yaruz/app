package session

import (
	"github.com/casdoor/casdoor-go-sdk/auth"
	"golang.org/x/oauth2"

	"github.com/yaruz/app/internal/domain/account"
	"github.com/yaruz/app/internal/domain/tg_account"
	"github.com/yaruz/app/internal/domain/user"
)

type Session struct {
	JwtClaims       *auth.Claims
	Token           *oauth2.Token
	User            *user.User
	AccountSettings *account.AccountSettings
	TgAccount       *tg_account.TgAccount
}

// New func is a constructor for the Session
func New() *Session {
	return &Session{}
}
