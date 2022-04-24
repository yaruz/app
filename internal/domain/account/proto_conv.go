package account

import (
	"github.com/casdoor/casdoor-go-sdk/auth"
	"github.com/yaruz/app/internal/pkg/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func AccountSettingsProto2AccountSettings(accountSettingsProto proto.AccountSettings) (accountSettings *AccountSettings, err error) {
	accountSettings = &AccountSettings{
		LangID: uint(accountSettingsProto.LangID),
	}
	return accountSettings, nil
}

func AccountSettings2AccountSettingsProto(accountSettings AccountSettings) (accountSettingsProto *proto.AccountSettings, err error) {
	accountSettingsProto = &proto.AccountSettings{
		LangID: uint64(accountSettings.LangID),
	}
	return accountSettingsProto, nil
}

func ClaimsProto2Claims(claimsProto proto.JwtClaims) (claims *auth.Claims, err error) {
	claims = &auth.Claims{
		AccessToken: claimsProto.AccessToken,
		Email:       claimsProto.Email,
	}
	if claimsProto.CreatedAt != nil && claimsProto.CreatedAt.IsValid() {
		claims.CreatedAt = claimsProto.CreatedAt.AsTime()
	}
	return claims, nil
}

func Claims2ClaimsProto(claims auth.Claims) (claimsProto *proto.JwtClaims, err error) {
	claimsProto = &proto.User{
		ID:        uint64(claims.ID),
		AccountID: claims.AccountID,
		Email:     claims.Email,
		CreatedAt: timestamppb.New(claims.CreatedAt),
	}
	return claimsProto, nil
}
