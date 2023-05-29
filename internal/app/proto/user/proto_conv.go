package user

import (
	"github.com/yaruz/app/internal/domain/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UserProto2User(userProto *User) (u *user.User) {
	if userProto == nil {
		return nil
	}

	u = &user.User{
		ID:        uint(userProto.ID),
		Email:     userProto.Email,
		Phone:     userProto.Phone,
		FirstName: userProto.FirstName,
		LastName:  userProto.LastName,
		UserName:  userProto.UserName,
	}
	if userProto.CreatedAt != nil && userProto.CreatedAt.IsValid() {
		u.CreatedAt = userProto.CreatedAt.AsTime()
	}
	return u
}

func User2UserProto(user *user.User) (up *User) {
	if user == nil {
		return nil
	}

	up = &User{
		ID:        uint64(user.ID),
		Email:     user.Email,
		Phone:     user.Phone,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserName:  user.UserName,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
	return up
}
