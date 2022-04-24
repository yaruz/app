package session

import (
	"github.com/casdoor/casdoor-go-sdk/auth"
	"github.com/yaruz/app/internal/domain/account"
	"github.com/yaruz/app/internal/domain/sn_account"

	"github.com/yaruz/app/internal/domain/user"
)

type Session struct {
	JwtClaims       *auth.Claims
	User            *user.User
	AccountSettings *account.AccountSettings
	SnAccounts      []sn_account.SNAccount
}
