package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/config-source/cdb/pkg/auth"
	"github.com/rs/zerolog"
)

type contextUserKey struct{}

const (
	IDTokenCookieName      = "cdb-id-token"
	AccessTokenCookieName  = "cdb-access-token"
	RefreshTokenCookieName = "cdb-session"

	AuthorizationHeaderPrefix = "Bearer "
)

func GetUser(r *http.Request) (auth.User, error) {
	user, ok := r.Context().Value(contextUserKey{}).(*auth.User)
	if !ok {
		return auth.User{}, errors.New("unable to get user from request context")
	}

	return *user, nil
}

func Authentication(log zerolog.Logger, userSvc *auth.UserService, signingKey []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var token string
			header := r.Header.Get("Authorization")
			if header != "" {
				if !strings.HasPrefix(header, AuthorizationHeaderPrefix) {
					w.WriteHeader(http.StatusBadRequest)
					// nolint:errcheck
					w.Write([]byte("Malformed Authorization header."))
					return
				}

				token = header[len(AuthorizationHeaderPrefix):]
			}

			if token == "" {
				cookie, err := r.Cookie(IDTokenCookieName)
				if err != nil && !errors.Is(err, http.ErrNoCookie) {
					w.WriteHeader(http.StatusBadRequest)
					// nolint:errcheck
					w.Write([]byte("unable to read session cookie"))
					log.Err(err).Msg("unable to read session cookie")
					return
				}

				if err == nil {
					token = cookie.Value
				} else {
					log.Debug().Str("cookieErr", err.Error()).Msg("error getting ID token cookie")
				}
			}

			if token == "" {
				log.Debug().Msg("unable to get token for request")
				next.ServeHTTP(w, r)
				return
			}

			user, err := auth.ValidateIdToken(signingKey, token)
			if err != nil {
				log.Err(err).Msg("invalid token")
				w.WriteHeader(http.StatusBadRequest)
				// nolint:errcheck
				w.Write([]byte("invalid token"))
				return
			}

			newCtx := context.WithValue(r.Context(), contextUserKey{}, &user)
			next.ServeHTTP(w, r.WithContext(newCtx))
		},
	)
}

func AuthenticationRequired(log zerolog.Logger, userSvc *auth.UserService, signingKey []byte, next http.Handler) http.Handler {
	return Authentication(
		log,
		userSvc,
		signingKey,
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if _, err := GetUser(r); err == nil {
					next.ServeHTTP(w, r)
					return
				}

				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Add("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(struct{ Message string }{Message: "forbidden"}); err != nil {
					log.Err(err).Msg("failed to encode a payload")
				}
			},
		),
	)
}
