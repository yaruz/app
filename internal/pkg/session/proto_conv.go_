package session

import (
	"github.com/golang/protobuf/ptypes"

	"github.com/yaruz/app/internal/pkg/proto"

	"github.com/yaruz/app/internal/domain/user"
)

func SessionProto2Session(sessionProto proto.Session) (s *Session, err error) {
	user, err := user.UserProto2User(*sessionProto.User)
	if err != nil {
		return nil, err
	}
	data, err := DataProto2Data(*sessionProto.Data)
	if err != nil {
		return nil, err
	}
	s = &Session{
		ID:     uint(sessionProto.ID),
		UserID: uint(sessionProto.UserID),
		Token:  sessionProto.Token,
		User:   *user,
		Data:   *data,
	}
	return s, nil
}

func Session2SessionProto(session Session) (sessionProto *proto.Session, err error) {
	userProto, err := user.User2UserProto(session.User)
	if err != nil {
		return nil, err
	}
	dataProto, err := Data2DataProto(session.Data)
	if err != nil {
		return nil, err
	}
	sessionProto = &proto.Session{
		ID:     uint64(session.ID),
		UserID: uint64(session.UserID),
		Token:  session.Token,
		User:   userProto,
		Data:   dataProto,
	}
	return sessionProto, nil
}

func DataProto2Data(dataProto proto.Data) (data *Data, err error) {
	data = &Data{
		UserID:   uint(dataProto.UserID),
		UserName: dataProto.UserName,
	}
	if dataProto.ExpirationTokenTime != nil {
		data.ExpirationTokenTime, err = ptypes.Timestamp(dataProto.ExpirationTokenTime)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func Data2DataProto(data Data) (dataProto *proto.Data, err error) {
	dataProto = &proto.Data{
		UserID:   uint64(data.UserID),
		UserName: data.UserName,
	}
	dataProto.ExpirationTokenTime, err = ptypes.TimestampProto(data.ExpirationTokenTime)
	if err != nil {
		return nil, err
	}
	return dataProto, nil
}
