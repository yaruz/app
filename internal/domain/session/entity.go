package session

import (
	"encoding"
	"github.com/casdoor/casdoor-go-sdk/auth"
	"github.com/yaruz/app/internal/domain/account"
	"github.com/yaruz/app/internal/domain/sn_account"
	"github.com/yaruz/app/internal/pkg/proto"
	"golang.org/x/oauth2"

	"github.com/yaruz/app/internal/domain/user"
)

type Session struct {
	JwtClaims       *auth.Claims
	Token           *oauth2.Token
	User            *user.User
	AccountSettings *account.AccountSettings
	SNAccounts      []sn_account.SNAccount
}

var _ encoding.BinaryMarshaler = (*Session)(nil)
var _ encoding.BinaryUnmarshaler = (*Session)(nil)

// New func is a constructor for the Session
func New() *Session {
	return &Session{}
}

func (e *Session) MarshalBinary() (data []byte, err error) {
	sessionProto, err := Session2SessionProto(e)
	if err != nil {
		return nil, err
	}
	return sessionProto.MarshalBinary()
}

func (e *Session) UnmarshalBinary(data []byte) (err error) {
	sessionProto := &proto.Session{}

	err = sessionProto.UnmarshalBinary(data)
	if err != nil {
		return err
	}

	s, err := SessionProto2Session(sessionProto)
	if err != nil {
		return err
	}

	*e = *s
	return nil
}
