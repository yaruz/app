package session

import (
	"github.com/yaruz/app/internal/domain/account"
	"github.com/yaruz/app/internal/domain/sn_account"
	"github.com/yaruz/app/internal/domain/user"
	"github.com/yaruz/app/internal/pkg/proto"
)

func SessionProto2Session(sessionProto *proto.Session) (s *Session, err error) {
	token, err := account.TokenProto2Token(sessionProto.Token)
	if err != nil {
		return nil, err
	}

	claims, err := account.ClaimsProto2Claims(sessionProto.JwtClaims)
	if err != nil {
		return nil, err
	}

	u, err := user.UserProto2User(sessionProto.User)
	if err != nil {
		return nil, err
	}

	accountSettings, err := account.AccountSettingsProto2AccountSettings(sessionProto.AccountSettings)
	if err != nil {
		return nil, err
	}

	sNAccounts := make([]sn_account.SNAccount, 0, len(sessionProto.SNAccounts))
	for _, snaccProto := range sessionProto.SNAccounts {
		snacc, err := sn_account.SNAccountProto2SNAccount(snaccProto)
		if err != nil {
			return nil, err
		}
		sNAccounts = append(sNAccounts, *snacc)
	}

	s = &Session{
		Token:           token,
		JwtClaims:       claims,
		User:            u,
		AccountSettings: accountSettings,
		SNAccounts:      sNAccounts,
	}
	return s, nil
}

func Session2SessionProto(s *Session) (sessionProto *proto.Session, err error) {
	tokenProto, err := account.Token2TokenProto(s.Token)
	if err != nil {
		return nil, err
	}

	claimsProto, err := account.Claims2ClaimsProto(s.JwtClaims)
	if err != nil {
		return nil, err
	}

	userProto, err := user.User2UserProto(s.User)
	if err != nil {
		return nil, err
	}

	accountSettingsProto, err := account.AccountSettings2AccountSettingsProto(s.AccountSettings)
	if err != nil {
		return nil, err
	}

	sNAccountsProto := make([]*proto.SNAccount, 0, len(s.SNAccounts))
	for _, snacc := range s.SNAccounts {
		snaccProto, err := sn_account.SNAccount2SNAccountProto(&snacc)
		if err != nil {
			return nil, err
		}
		sNAccountsProto = append(sNAccountsProto, snaccProto)
	}

	sessionProto = &proto.Session{
		User:            userProto,
		AccountSettings: accountSettingsProto,
		SNAccounts:      sNAccountsProto,
		JwtClaims:       claimsProto,
		Token:           tokenProto,
	}
	return sessionProto, nil
}
