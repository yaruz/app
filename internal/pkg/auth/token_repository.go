package auth

import (
	"time"
)

type TokenRepository interface {
	NewTokenByData(data TokenData) Token
	ParseStringToken(tokenString string, signingKey string) (Token, error)
}

type Token interface {
	GetData() TokenData
	GenerateStringToken(signingKey string) (string, error)
	Valid() error
}

type TokenData struct {
	UserID              uint
	UserName            string
	ExpirationTokenTime time.Time
}
