package session

import (
	account_proto "github.com/yaruz/app/internal/app/proto/account"
	jwt_proto "github.com/yaruz/app/internal/app/proto/jwt"
	user_proto "github.com/yaruz/app/internal/app/proto/user"
	
	"github.com/yaruz/app/internal/pkg/session"
)

func SessionProto2Session(sessionProto *Session) (s *session.Session) {
	if sessionProto == nil {
		return nil
	}

	return &session.Session{
		ID:              sessionProto.ID,
		AccountSettings: account_proto.AccountSettingsProto2AccountSettings(sessionProto.AccountSettings),
		JwtClaims:       jwt_proto.ClaimsProto2Claims(sessionProto.JwtClaims),
		User:            user_proto.UserProto2User(sessionProto.User),
	}
}

func Session2SessionProto(s *session.Session) (sessionProto *Session) {
	if s == nil {
		return nil
	}

	return &Session{
		ID:              s.ID,
		AccountSettings: account_proto.AccountSettings2AccountSettingsProto(s.AccountSettings),
		JwtClaims:       jwt_proto.Claims2ClaimsProto(s.JwtClaims),
		User:            user_proto.User2UserProto(s.User),
	}
}
