package tg_account

import (
	mtproto_session "github.com/Kalinin-Andrey/mtproto/session"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/yaruz/app/internal/domain/tg_account"
)

func TgAccountProto2TgAccount(tgAccountProto *TgAccount) (tgAccount *tg_account.TgAccount, err error) {
	if tgAccountProto == nil {
		return nil, nil
	}

	authSession, err := AuthSessionProto2AuthSession(tgAccountProto.AuthSession)
	if err != nil {
		return nil, err
	}
	tgAccount = &tg_account.TgAccount{
		ID:          uint(tgAccountProto.ID),
		TgID:        tgAccountProto.TgID,
		AuthSession: authSession,
	}
	if tgAccountProto.CreatedAt != nil && tgAccountProto.CreatedAt.IsValid() {
		tgAccount.CreatedAt = tgAccountProto.CreatedAt.AsTime()
	}
	return tgAccount, nil
}

func TgAccount2TgAccountProto(tgAccount *tg_account.TgAccount) (tgAccountProto *TgAccount, err error) {
	if tgAccount == nil {
		return nil, nil
	}

	authSessionProto, err := AuthSession2AuthSessionProto(tgAccount.AuthSession)
	if err != nil {
		return nil, err
	}
	tgAccountProto = &TgAccount{
		ID:          uint64(tgAccount.ID),
		TgID:        tgAccount.TgID,
		AuthSession: authSessionProto,
		CreatedAt:   timestamppb.New(tgAccount.CreatedAt),
	}
	return tgAccountProto, nil
}

func AuthSessionProto2AuthSession(authSessionProto *AuthSession) (authSession *mtproto_session.Session, err error) {
	if authSessionProto == nil {
		return nil, nil
	}

	authSession = &mtproto_session.Session{
		Key:      authSessionProto.Key,
		Hash:     authSessionProto.Hash,
		Salt:     authSessionProto.Salt,
		Hostname: authSessionProto.Hostname,
	}
	return authSession, nil
}

func AuthSession2AuthSessionProto(authSession *mtproto_session.Session) (authSessionProto *AuthSession, err error) {
	if authSession == nil {
		return nil, nil
	}

	authSessionProto = &AuthSession{
		Key:      authSession.Key,
		Hash:     authSession.Hash,
		Salt:     authSession.Salt,
		Hostname: authSession.Hostname,
	}
	return authSessionProto, nil
}
