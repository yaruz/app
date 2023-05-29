package auth

import (
	"github.com/yaruz/app/internal/pkg/jwt"
)

type Token interface {
	GetClaims() *jwt.Claims
	GetData() *jwt.TokenData
	GenerateStringToken() (string, error)
}

type TokenRepository interface {
	NewTokenWithData(data *jwt.TokenData) (*jwt.Token, error)
	ParseStringToken(tokenString string) (*jwt.Token, error)
}

var _ TokenRepository = (*jwt.Repository)(nil)
var _ Token = (*jwt.Token)(nil)
