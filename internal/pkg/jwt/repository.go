package jwt

import (
	"github.com/yaruz/app/internal/domain/user"
	"strconv"
	"time"

	"github.com/pkg/errors"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/yaruz/app/internal/pkg/apperror"
)

type Token struct {
	signingKey string
	claims     Claims
}

func (t Token) GenerateStringToken() (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, t.claims).SignedString([]byte(t.signingKey))
}

func (t Token) GetData() *TokenData {
	return t.claims.TokenData
}

func (t Token) GetClaims() *Claims {
	return &t.claims
}

type Claims struct {
	*jwt.RegisteredClaims
	*TokenData
	AccessToken string
}

type TokenData struct {
	SessionID string
	User      *user.User
}

type Repository struct {
	signingKey        string
	expirationInHours uint
}

func NewRepository(signingKey string, expirationInHours uint) *Repository {
	return &Repository{
		signingKey:        signingKey,
		expirationInHours: expirationInHours,
	}
}

func (r Repository) NewTokenWithData(data *TokenData) (*Token, error) {
	var err error
	token := &Token{
		signingKey: r.signingKey,
		claims: Claims{
			TokenData: data,
			RegisteredClaims: &jwt.RegisteredClaims{
				Subject: strconv.Itoa(int(data.User.ID)),
				ExpiresAt: &jwt.NumericDate{
					time.Now().Add(time.Hour * time.Duration(r.expirationInHours)).UTC(),
				},
			},
		},
	}
	token.claims.AccessToken, err = token.GenerateStringToken()
	return token, err
}

func (r Repository) ParseStringToken(accessToken string) (*Token, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(r.signingKey), nil
	})
	if err != nil {
		return nil, errors.Wrapf(apperror.ErrInvalidToken, "Text.ParseToken error: %v", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.Wrapf(apperror.ErrInvalidToken, "Text.ParseToken error.")
	}
	claims.AccessToken = accessToken

	if !token.Valid {
		return nil, apperror.ErrTokenHasExpired
	}

	return &Token{
		claims:     *claims,
		signingKey: r.signingKey,
	}, nil
}

//  generates a JWT that encodes an identity.
/*func (r sessionRepository) generateJWT(user *user.User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"name": user.Name,
		"exp":  time.Now().Add(time.Duration(r.tokenExpiration) * time.Hour).Unix(),
	}).SignedString([]byte(r.signingKey))
}*/
