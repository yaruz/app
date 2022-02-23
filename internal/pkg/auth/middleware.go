package auth

import (
	"github.com/yaruz/app/internal/domain/user"
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
		token := ""
		var err error

		if strings.HasPrefix(header, "Bearer ") {
			token = header[7:]
		} else {
			return UnauthorizedError(rctx, "")
		}

		if err = authService.StringTokenValidation(ctx, token); err != nil {
			return UnauthorizedError(rctx, err.Error())
		}

		if ctx, err = authService.SessionInit(ctx, token, getAccountSettings(rctx)); err != nil {
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

func getAccountSettings(rctx *routing.Context) *user.AccountSettings {
	// todo
	return nil
}

// CurrentUser returns the user identity from the given context.
// Nil is returned if no user identity is found in the context.
//func CurrentSession(ctx context.Context) *session.Session {
//	if sess, ok := ctx.Value(userSessionKey).(*session.Session); ok {
//		return sess
//	}
//	return nil
//}
