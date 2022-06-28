package user

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/yaruz/app/internal/domain/user"
)

func UserProto2User(userProto *User) (u *user.User, err error) {
	u = &user.User{
		ID:        uint(userProto.ID),
		AccountID: userProto.AccountID,
		Email:     userProto.Email,
	}
	if userProto.CreatedAt != nil && userProto.CreatedAt.IsValid() {
		u.CreatedAt = userProto.CreatedAt.AsTime()
	}
	return u, nil
}

func User2UserProto(user *user.User) (up *User, err error) {
	up = &User{
		ID:        uint64(user.ID),
		AccountID: user.AccountID,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
	return up, nil
}
