package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/config-source/cdb/internal/auth"
	"github.com/rs/zerolog"
)

const SessionCookieName = "cdb-session"
const AuthorizationHeaderPrefix = "Bearer "
const contextUserKey = "user"

func GetUser(r *http.Request) auth.User {
	user, ok := r.Context().Value(contextUserKey).(auth.User)
	if !ok {
		panic("somehow reached unreachable code in GetUser!")
	}

	return user
}

// TODO: tests

func Authentication(log zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var token string
			header := r.Header.Get("Authorization")
			if header != "" {
				if !strings.HasPrefix(header, AuthorizationHeaderPrefix) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("Malformed Authorization header."))
					return
				}

				token = header[len(AuthorizationHeaderPrefix):]
			}

			if token == "" {
				cookie, err := r.Cookie(SessionCookieName)
				if err != nil && !errors.Is(err, http.ErrNoCookie) {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("unable to read session cookie"))
					log.Err(err).Msg("unable to read session cookie")
					return
				}

				if err == nil {
					token = cookie.Value
				}
			}

			user, err := auth.ValidateIdToken(token)
			if err != nil {
				log.Err(err).Msg("invalid token")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("invalid token"))
				return
			}

			newCtx := context.WithValue(r.Context(), contextUserKey, user)
			next.ServeHTTP(w, r.WithContext(newCtx))
		},
	)
}
