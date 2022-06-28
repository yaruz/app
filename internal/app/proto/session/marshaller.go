package session

import "github.com/yaruz/app/internal/domain/session"

func MarshalBinary(e *session.Session) ([]byte, error) {
	sessionProto, err := Session2SessionProto(e)
	if err != nil {
		return nil, err
	}
	return sessionProto.MarshalBinary()
}

func UnmarshalBinary(data []byte) (*session.Session, error) {
	sessionProto := &Session{}

	err := sessionProto.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return SessionProto2Session(sessionProto)
}
