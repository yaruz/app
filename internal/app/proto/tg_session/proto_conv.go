package tg_session

import (
	"github.com/yaruz/app/internal/pkg/socnets/tg"
)

func TgSessionProto2TgSession(tgSessionProto *TgSession) (tgSession *tg.Session) {
	if tgSessionProto == nil {
		return nil
	}

	return &tg.Session{
		IsAuthorized:  tgSessionProto.IsAuthorized,
		Session:       tgSessionProto.Session,
		ID:            tgSessionProto.ID,
		PhoneCodeHash: tgSessionProto.PhoneCodeHash,
		Phone:         tgSessionProto.Phone,
	}
}

func TgSession2TgSessionProto(tgSession *tg.Session) (tgSessionProto *TgSession) {
	if tgSession == nil {
		return nil
	}

	return &TgSession{
		IsAuthorized:  tgSession.IsAuthorized,
		Session:       tgSession.Session,
		ID:            tgSession.ID,
		Phone:         tgSession.Phone,
		PhoneCodeHash: tgSession.PhoneCodeHash,
	}
}
