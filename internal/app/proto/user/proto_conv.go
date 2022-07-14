package user

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/yaruz/app/internal/domain/user"
)

func UserProto2User(userProto *User) (u *user.User, err error) {
	if userProto == nil {
		return nil, nil
	}

	u = &user.User{
		ID:        uint(userProto.ID),
		AccountID: userProto.AccountID,
		Email:     userProto.Email,
		Phone:     userProto.Phone,
	}
	if userProto.CreatedAt != nil && userProto.CreatedAt.IsValid() {
		u.CreatedAt = userProto.CreatedAt.AsTime()
	}
	return u, nil
}

func User2UserProto(user *user.User) (up *User, err error) {
	if user == nil {
		return nil, nil
	}

	up = &User{
		ID:        uint64(user.ID),
		AccountID: user.AccountID,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
	return up, nil
}
