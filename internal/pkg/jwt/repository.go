package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"strconv"

	"github.com/yaruz/app/internal/pkg/apperror"
	"github.com/yaruz/app/internal/pkg/auth"
)

type Repository struct {
}

type Token struct {
	Data   auth.TokenData
	claims claims
}

type claims struct {
	jwt.StandardClaims
	auth.TokenData
}

var _ auth.TokenRepository = (*Repository)(nil)
var _ auth.Token = (*Token)(nil)

func NewRepository() *Repository {
	return &Repository{}
}

func (r Repository) NewTokenByData(data auth.TokenData) auth.Token {
	return &Token{
		Data: data,
		claims: claims{
			TokenData: data,
			StandardClaims: jwt.StandardClaims{
				Subject:   strconv.Itoa(int(data.UserID)),
				ExpiresAt: data.ExpirationTokenTime.Unix(),
			},
		},
	}
}

func (r Repository) ParseStringToken(tokenString string, signingKey string) (auth.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		return nil, errors.Wrapf(apperror.ErrBadRequest, "Text.ParseToken error: %v", err)
	}

	claims, ok := token.Claims.(*claims)
	if !ok {
		return nil, errors.Wrapf(apperror.ErrBadRequest, "Text.ParseToken error.")
	}

	if !token.Valid {
		return nil, apperror.ErrTokenHasExpired
	}

	return &Token{
		claims: *claims,
		Data: auth.TokenData{
			UserID:              claims.UserID,
			UserName:            claims.UserName,
			ExpirationTokenTime: claims.ExpirationTokenTime,
		},
	}, nil
}

func (t Token) GenerateStringToken(signingKey string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, t.claims).SignedString([]byte(signingKey))
}

func (t Token) GetData() auth.TokenData {
	return t.Data
}

func (t Token) Valid() error {
	return t.claims.Valid()
}

//  generates a JWT that encodes an identity.
/*func (r SessionRepository) generateJWT(user *user.User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"name": user.Name,
		"exp":  time.Now().Add(time.Duration(r.tokenExpiration) * time.Hour).Unix(),
	}).SignedString([]byte(r.signingKey))
}*/
