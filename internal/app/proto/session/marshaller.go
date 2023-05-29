package session

import (
	"github.com/yaruz/app/internal/pkg/session"
)

func MarshalBinary(e *session.Session) ([]byte, error) {
	return Session2SessionProto(e).MarshalBinary()
}

func UnmarshalBinary(data []byte) (*session.Session, error) {
	sessionProto := &Session{}

	err := sessionProto.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return SessionProto2Session(sessionProto), nil
}
