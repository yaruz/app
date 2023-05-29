package jwt

import (
	gojwt "github.com/golang-jwt/jwt/v5"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/yaruz/app/internal/pkg/jwt"

	user_proto "github.com/yaruz/app/internal/app/proto/user"
)

func ClaimsProto2Claims(claimsProto *Claims) (claims *jwt.Claims) {
	if claimsProto == nil {
		return nil
	}

	return &jwt.Claims{
		RegisteredClaims: RegisteredClaimsProto2RegisteredClaims(claimsProto.RegisteredClaims),
		TokenData:        TokenDataProto2TokenData(claimsProto.TokenData),
		AccessToken:      claimsProto.AccessToken,
	}
}
func Claims2ClaimsProto(claims *jwt.Claims) (claimsProto *Claims) {
	if claims == nil {
		return nil
	}

	return &Claims{
		RegisteredClaims: RegisteredClaims2RegisteredClaimsProto(claims.RegisteredClaims),
		TokenData:        TokenData2TokenDataProto(claims.TokenData),
		AccessToken:      claims.AccessToken,
	}
}

func RegisteredClaimsProto2RegisteredClaims(registeredClaimsProto *RegisteredClaims) (registeredClaims *gojwt.RegisteredClaims) {
	if registeredClaimsProto == nil {
		return nil
	}

	registeredClaims = &gojwt.RegisteredClaims{
		ID:       registeredClaimsProto.ID,
		Issuer:   registeredClaimsProto.Issuer,
		Subject:  registeredClaimsProto.Subject,
		Audience: gojwt.ClaimStrings(registeredClaimsProto.Audience),
	}
	if registeredClaimsProto.ExpiresAt != nil && registeredClaimsProto.ExpiresAt.IsValid() {
		registeredClaims.ExpiresAt = &gojwt.NumericDate{Time: registeredClaimsProto.ExpiresAt.AsTime()}
	}
	if registeredClaimsProto.NotBefore != nil && registeredClaimsProto.NotBefore.IsValid() {
		registeredClaims.NotBefore = &gojwt.NumericDate{Time: registeredClaimsProto.NotBefore.AsTime()}
	}
	if registeredClaimsProto.IssuedAt != nil && registeredClaimsProto.IssuedAt.IsValid() {
		registeredClaims.IssuedAt = &gojwt.NumericDate{Time: registeredClaimsProto.IssuedAt.AsTime()}
	}
	return registeredClaims
}
func RegisteredClaims2RegisteredClaimsProto(registeredClaims *gojwt.RegisteredClaims) (registeredClaimsProto *RegisteredClaims) {
	if registeredClaims == nil {
		return nil
	}

	return &RegisteredClaims{
		ID:        registeredClaims.ID,
		Issuer:    registeredClaims.Issuer,
		Subject:   registeredClaims.Subject,
		Audience:  registeredClaims.Audience,
		ExpiresAt: timestamppb.New(registeredClaims.ExpiresAt.Time),
		NotBefore: timestamppb.New(registeredClaims.NotBefore.Time),
		IssuedAt:  timestamppb.New(registeredClaims.IssuedAt.Time),
	}
}

func TokenDataProto2TokenData(tokenDataProto *TokenData) (tokenData *jwt.TokenData) {
	if tokenDataProto == nil {
		return nil
	}

	return &jwt.TokenData{
		SessionID: tokenDataProto.SessionID,
		User:      user_proto.UserProto2User(tokenDataProto.User),
	}
}
func TokenData2TokenDataProto(tokenData *jwt.TokenData) (tokenDataProto *TokenData) {
	if tokenData == nil {
		return nil
	}

	return &TokenData{
		SessionID: tokenData.SessionID,
		User:      user_proto.User2UserProto(tokenData.User),
	}
}
