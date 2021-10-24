package user

import (
	"github.com/golang/protobuf/ptypes"

	"github.com/yaruz/app/internal/pkg/proto"
)

func UserProto2User(userProto proto.User) (u *User, err error) {
	u = &User{
		ID:   uint(userProto.ID),
		Name: userProto.Name,
	}
	if userProto.CreatedAt != nil {
		u.CreatedAt, err = ptypes.Timestamp(userProto.CreatedAt)
		if err != nil {
			return nil, err
		}
	}
	if userProto.UpdatedAt != nil {
		u.UpdatedAt, err = ptypes.Timestamp(userProto.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}
	if userProto.DeletedAt != nil {
		t, err := ptypes.Timestamp(userProto.DeletedAt)
		if err != nil {
			return nil, err
		}
		u.DeletedAt = &t
	}
	return u, nil
}

func User2UserProto(user User) (up *proto.User, err error) {
	up = &proto.User{
		ID:   uint64(user.ID),
		Name: user.Name,
	}
	up.CreatedAt, err = ptypes.TimestampProto(user.CreatedAt)
	if err != nil {
		return nil, err
	}
	up.UpdatedAt, err = ptypes.TimestampProto(user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if user.DeletedAt != nil {
		up.DeletedAt, err = ptypes.TimestampProto(*user.DeletedAt)
		if err != nil {
			return nil, err
		}
	}
	return up, nil
}
