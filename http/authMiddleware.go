package http

import (
	"context"
	"net/http"

	api "github.com/eriktate/splitdecision"
)

type sessionKey struct{}

func WithSession(r *http.Request, session api.Session) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), sessionKey{}, session))
}

func GetSessionFromContext(ctx context.Context) (api.Session, bool) {
	possibleSession := ctx.Value(sessionKey{})
	if possibleSession != nil {
		session, ok := possibleSession.(api.Session)
		return session, ok
	}

	return api.Session{}, false
}

func RequireSession(sessionService api.SessionService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			cookie, err := r.Cookie("session")
			if err != nil {
				unauthorized(w, "no valid session, please login")
				return
			}

			sessionID, err := api.ParseID(cookie.Value)
			if err != nil {
				unauthorized(w, "no valid session, please login")
				return
			}

			session, err := sessionService.GetSession(ctx, sessionID)
			if err != nil {
				unauthorized(w, "no valid session, please login")
				return
			}

			if !session.Valid() {
				unauthorized(w, "no valid session, please login")
				return
			}

			next.ServeHTTP(w, WithSession(r, session))
		})
	}
}

func IncludeSession(sessionService api.SessionService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer next.ServeHTTP(w, r)
			ctx := r.Context()
			cookie, err := r.Cookie("session")
			if err != nil {
				return
			}

			sessionID, err := api.ParseID(cookie.Value)
			if err != nil {
				return
			}

			session, err := sessionService.GetSession(ctx, sessionID)
			if err != nil {
				return
			}

			if !session.Valid() {
				return
			}

			r = WithSession(r, session)
		})
	}
}
