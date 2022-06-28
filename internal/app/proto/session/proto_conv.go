package session

import (
	account_proto "github.com/yaruz/app/internal/app/proto/account"
	snaccount_proto "github.com/yaruz/app/internal/app/proto/sn_account"
	token_proto "github.com/yaruz/app/internal/app/proto/token"
	user_proto "github.com/yaruz/app/internal/app/proto/user"

	"github.com/yaruz/app/internal/domain/session"
	"github.com/yaruz/app/internal/domain/sn_account"
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

	sNAccounts := make([]sn_account.SNAccount, 0, len(sessionProto.SNAccounts))
	for _, snaccProto := range sessionProto.SNAccounts {
		snacc, err := snaccount_proto.SNAccountProto2SNAccount(snaccProto)
		if err != nil {
			return nil, err
		}
		sNAccounts = append(sNAccounts, *snacc)
	}

	s = &session.Session{
		Token:           token,
		JwtClaims:       claims,
		User:            u,
		AccountSettings: accountSettings,
		SNAccounts:      sNAccounts,
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

	sNAccountsProto := make([]*snaccount_proto.SNAccount, 0, len(s.SNAccounts))
	for _, snacc := range s.SNAccounts {
		snaccProto, err := snaccount_proto.SNAccount2SNAccountProto(&snacc)
		if err != nil {
			return nil, err
		}
		sNAccountsProto = append(sNAccountsProto, snaccProto)
	}

	sessionProto = &Session{
		User:            userProto,
		AccountSettings: accountSettingsProto,
		SNAccounts:      sNAccountsProto,
		JwtClaims:       claimsProto,
		Token:           tokenProto,
	}
	return sessionProto, nil
}
