package session

import (
	"github.com/casdoor/casdoor-go-sdk/auth"
)

type Session struct {
	JwtClaims auth.Claims
}
