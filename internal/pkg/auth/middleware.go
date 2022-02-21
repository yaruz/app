package auth

import (
	"net/http"
	"strings"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/minipkg/log"
)

// Middleware returns a JWT-based authentication middleware.
func Middleware(logger log.ILogger, authService Service) routing.Handler {
	return func(rctx *routing.Context) error {
		ctx := rctx.Request.Context()
		header := rctx.Request.Header.Get("Authorization")
		message := ""
		if strings.HasPrefix(header, "Bearer ") {
			ctx, ok, err := authService.StringTokenValidation(ctx, header[7:])
			if err == nil && ok {
				*rctx.Request = *rctx.Request.WithContext(ctx)
				return nil
			}
			if err != nil {
				message = err.Error()
			}

		}
		rctx.Response.Header().Set("WWW-Authenticate", `Bearer realm="API"`)
		if message != "" {
			return routing.NewHTTPError(http.StatusUnauthorized, message)
		}
		return routing.NewHTTPError(http.StatusUnauthorized)
	}
}

// CurrentUser returns the user identity from the given context.
// Nil is returned if no user identity is found in the context.
//func CurrentSession(ctx context.Context) *session.Session {
//	if sess, ok := ctx.Value(userSessionKey).(*session.Session); ok {
//		return sess
//	}
//	return nil
//}
