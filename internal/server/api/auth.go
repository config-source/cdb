package api

import (
	"encoding/json"
	"net/http"

	"github.com/config-source/cdb/internal/auth"
	"github.com/config-source/cdb/internal/server/middleware"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *API) doLogin(w http.ResponseWriter, r *http.Request, user auth.User) {
	tokens, err := auth.GenerateTokens(a.tokenSigningKey, user)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	authCookies := map[string]string{
		middleware.IDTokenCookieName:      tokens.IDToken,
		middleware.AccessTokenCookieName:  tokens.AccessToken,
		middleware.RefreshTokenCookieName: tokens.RefreshToken,
	}

	for name, value := range authCookies {
		http.SetCookie(
			w,
			&http.Cookie{
				Name:     name,
				Value:    value,
				Domain:   r.Host,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			},
		)
	}

	a.sendJson(w, tokens)
}

func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var request Credentials
	err := decoder.Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, err)
		return
	}

	user, err := a.userService.Login(r.Context(), request.Email, request.Password)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.doLogin(w, r, user)
}

func (a *API) Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var request Credentials
	err := decoder.Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, err)
		return
	}

	user, err := a.userService.Register(r.Context(), request.Email, request.Password)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.doLogin(w, r, user)
}
