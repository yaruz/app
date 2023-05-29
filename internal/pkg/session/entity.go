package session

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/yaruz/app/internal/domain/account"
	"github.com/yaruz/app/internal/domain/user"
	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/internal/pkg/jwt"
)

type Session struct {
	ID              string                   // обязательное
	AccountSettings *account.AccountSettings // обязательное
	JwtClaims       *jwt.Claims              // только после первой аутентификации (но факт аутентификации проверяем по переменной в контексте isAuth)
	User            *user.User               // только после первой аутентификации (но факт аутентификации проверяем по переменной в контексте isAuth)
}

// New func is a constructor for the Session
func New(ID string, accountSettings *account.AccountSettings, jwtClaims *jwt.Claims, user *user.User) (*Session, error) {
	var err error
	if ID == "" {
		ID, err = generateSessionID()
	} else {
		_, err = uuid.Parse(ID)
		err = fmt.Errorf("[%w] uuid.Parse error: %s", apperror.ErrBadParams, err.Error())
	}
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:              ID,
		AccountSettings: accountSettings,
		JwtClaims:       jwtClaims,
		User:            user,
	}, nil
}

func generateSessionID() (string, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("[%w] uuid.NewUUID error: %s", apperror.ErrInternal, err.Error())
	}
	return uuid.String(), nil
}
