package user

import (
	"github.com/yaruz/app/internal/pkg/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UserProto2User(userProto proto.User) (u *User, err error) {
	u = &User{
		ID:        uint(userProto.ID),
		AccountID: userProto.AccountID,
		Email:     userProto.Email,
	}
	if userProto.CreatedAt != nil && userProto.CreatedAt.IsValid() {
		u.CreatedAt = userProto.CreatedAt.AsTime()
	}
	return u, nil
}

func User2UserProto(user User) (up *proto.User, err error) {
	up = &proto.User{
		ID:        uint64(user.ID),
		AccountID: user.AccountID,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
	return up, nil
}
