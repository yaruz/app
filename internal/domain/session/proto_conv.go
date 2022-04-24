package session

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/yaruz/app/internal/pkg/proto"
)

func SessionProto2Session(sessionProto proto.Session) (s *Session, err error) {
	s = &Session{
		ID:   uint(userProto.ID),
		Name: userProto.Name,
	}
	if userProto.CreatedAt != nil {
		s.CreatedAt, err = ptypes.Timestamp(userProto.CreatedAt)
		if err != nil {
			return nil, err
		}
	}
	if userProto.UpdatedAt != nil {
		s.UpdatedAt, err = ptypes.Timestamp(userProto.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}
	if userProto.DeletedAt != nil {
		t, err := ptypes.Timestamp(userProto.DeletedAt)
		if err != nil {
			return nil, err
		}
		s.DeletedAt = &t
	}
	return s, nil
}

func Session2SessionProto(s Session) (sessionProto *proto.Session, err error) {
	sessionProto = &proto.User{
		ID:   uint64(s.ID),
		Name: s.Name,
	}
	sessionProto.CreatedAt, err = ptypes.TimestampProto(s.CreatedAt)
	if err != nil {
		return nil, err
	}
	sessionProto.UpdatedAt, err = ptypes.TimestampProto(s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if s.DeletedAt != nil {
		sessionProto.DeletedAt, err = ptypes.TimestampProto(*s.DeletedAt)
		if err != nil {
			return nil, err
		}
	}
	return sessionProto, nil
}
