package auth

import (
	"github.com/yaruz/app/internal/pkg/session"
	"context"
	"net/http"
	"strings"

	"github.com/minipkg/log"

	routing "github.com/go-ozzo/ozzo-routing/v2"
)

// Middleware returns a JWT-based authentication middleware.
func Middleware(logger log.ILogger, authService Service) routing.Handler {
	return func(c *routing.Context) error {
		ctx := c.Request.Context()
		header := c.Request.Header.Get("Authorization")
		message := ""
		if strings.HasPrefix(header, "Bearer ") {
			ctx, ok, err := authService.StringTokenValidation(ctx, header[7:])
			if err == nil && ok {
				*c.Request = *c.Request.WithContext(ctx)
				return nil
			}
			if err != nil {
				message = err.Error()
			}

		}
		c.Response.Header().Set("WWW-Authenticate", `Bearer realm="API"`)
		if message != "" {
			return routing.NewHTTPError(http.StatusUnauthorized, message)
		}
		return routing.NewHTTPError(http.StatusUnauthorized)
	}
}

// CurrentUser returns the user identity from the given context.
// Nil is returned if no user identity is found in the context.
func CurrentSession(ctx context.Context) *session.Session {
	if sess, ok := ctx.Value(userSessionKey).(*session.Session); ok {
		return sess
	}
	return nil
}
