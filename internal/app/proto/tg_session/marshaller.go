package tg_session

import (
	"github.com/yaruz/app/internal/pkg/socnets/tg"
)

func MarshalBinary(e *tg.Session) ([]byte, error) {
	return TgSession2TgSessionProto(e).MarshalBinary()
}

func UnmarshalBinary(data []byte) (*tg.Session, error) {
	tgSessionProto := &TgSession{}

	err := tgSessionProto.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return TgSessionProto2TgSession(tgSessionProto), nil
}
