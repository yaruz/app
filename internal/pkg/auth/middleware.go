package auth

import (
	"net/http"
	"strings"

	routing "github.com/go-ozzo/ozzo-routing/v2"
)

const (
	header_Authorization = "Authorization"
	header_Cookie        = "Cookie"
	cookie_SessionID     = "SessionID"
)

func getAuthorizationToken(headerAuthorization string) (token string) {
	if strings.HasPrefix(headerAuthorization, "Bearer ") {
		token = headerAuthorization[7:]
	}
	return token
}

func getSessionID(cookie string) string {
	return getCookie(cookie, cookie_SessionID)
}

func getCookie(cookie string, name string) string {
	var st string
	var p []string

	s := strings.Split(cookie, ";")
	m := make(map[string]string, len(s))
	for _, st = range s {
		p = strings.Split(st, "=")
		if len(p) != 2 {
			continue
		}
		m[strings.TrimSpace(p[0])] = strings.TrimSpace(p[1])
	}

	return m[name]
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
//func CurrentSession(ctx context.Context) *sessionRepository.Session {
//	if sess, ok := ctx.Value(ctxUser_SessionKey).(*sessionRepository.Session); ok {
//		return sess
//	}
//	return nil
//}
