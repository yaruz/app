package auth

import (
	"net/http"
	"strings"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/minipkg/log"
)

// Middleware returns a JWT-based authentication middleware.
func Middleware(logger log.ILogger, s Service) routing.Handler {
	return func(rctx *routing.Context) error {
		// todo: все сообщения при ошибках нужно выдавать не наружу, а в логи
		ctx := rctx.Request.Context()
		header := rctx.Request.Header.Get("Authorization")
		token := ""
		var err error

		if strings.HasPrefix(header, "Bearer ") {
			token = header[7:]
		} else {
			return UnauthorizedError(rctx, "")
		}

		if err = s.StringTokenValidation(ctx, token); err != nil {
			return UnauthorizedError(rctx, err.Error())
		}

		accountSettings, err := s.RoutingGetAccountSettingsWithDefaults(rctx)
		if err != nil {
			return routing.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if ctx, err = s.SessionInit(ctx, token, accountSettings); err != nil {
			return UnauthorizedError(rctx, err.Error())
		}

		*rctx.Request = *rctx.Request.WithContext(ctx)
		return nil
	}
}

func UnauthorizedError(rctx *routing.Context, message string) routing.HTTPError {
	rctx.Response.Header().Set("WWW-Authenticate", `Bearer realm="API"`)

	if message != "" {
		return routing.NewHTTPError(http.StatusUnauthorized, message)
	}
	return routing.NewHTTPError(http.StatusUnauthorized)
}

// CurrentUser returns the user identity from the given context.
// Nil is returned if no user identity is found in the context.
//func CurrentSession(ctx context.Context) *session.Session {
//	if sess, ok := ctx.Value(userSessionKey).(*session.Session); ok {
//		return sess
//	}
//	return nil
//}
