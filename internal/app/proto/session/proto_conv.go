package session

import (
	account_proto "github.com/yaruz/app/internal/app/proto/account"
	tg_account_proto "github.com/yaruz/app/internal/app/proto/tg_account"
	token_proto "github.com/yaruz/app/internal/app/proto/token"
	user_proto "github.com/yaruz/app/internal/app/proto/user"

	"github.com/yaruz/app/internal/domain/session"
)

func SessionProto2Session(sessionProto *Session) (s *session.Session, err error) {
	token, err := token_proto.TokenProto2Token(sessionProto.Token)
	if err != nil {
		return nil, err
	}

	claims, err := account_proto.ClaimsProto2Claims(sessionProto.JwtClaims)
	if err != nil {
		return nil, err
	}

	u, err := user_proto.UserProto2User(sessionProto.User)
	if err != nil {
		return nil, err
	}

	accountSettings, err := account_proto.AccountSettingsProto2AccountSettings(sessionProto.AccountSettings)
	if err != nil {
		return nil, err
	}

	tgAccount, err := tg_account_proto.TgAccountProto2TgAccount(sessionProto.TgAccount)
	if err != nil {
		return nil, err
	}

	s = &session.Session{
		Token:           token,
		JwtClaims:       claims,
		User:            u,
		AccountSettings: accountSettings,
		TgAccount:       tgAccount,
	}
	return s, nil
}

func Session2SessionProto(s *session.Session) (sessionProto *Session, err error) {
	tokenProto, err := token_proto.Token2TokenProto(s.Token)
	if err != nil {
		return nil, err
	}

	claimsProto, err := account_proto.Claims2ClaimsProto(s.JwtClaims)
	if err != nil {
		return nil, err
	}

	userProto, err := user_proto.User2UserProto(s.User)
	if err != nil {
		return nil, err
	}

	accountSettingsProto, err := account_proto.AccountSettings2AccountSettingsProto(s.AccountSettings)
	if err != nil {
		return nil, err
	}

	tgAccountsProto, err := tg_account_proto.TgAccount2TgAccountProto(s.TgAccount)
	if err != nil {
		return nil, err
	}

	sessionProto = &Session{
		User:            userProto,
		AccountSettings: accountSettingsProto,
		TgAccount:       tgAccountsProto,
		JwtClaims:       claimsProto,
		Token:           tokenProto,
	}
	return sessionProto, nil
}
