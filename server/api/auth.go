package api

import (
	"encoding/json"
	"net/http"

	"github.com/config-source/cdb/auth"
	"github.com/config-source/cdb/server/middleware"
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

	var creds Credentials
	err := decoder.Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, err)
		return
	}

	user, err := a.userService.Login(r.Context(), creds.Email, creds.Password)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.doLogin(w, r, user)
}

func (a *API) Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var creds Credentials
	err := decoder.Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, err)
		return
	}

	user, err := a.userService.Register(r.Context(), creds.Email, creds.Password)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.doLogin(w, r, user)
}

func (a *API) Logout(w http.ResponseWriter, r *http.Request) {
	authCookies := []string{
		middleware.IDTokenCookieName,
		middleware.AccessTokenCookieName,
		middleware.RefreshTokenCookieName,
	}

	for _, cookie := range authCookies {
		http.SetCookie(
			w,
			&http.Cookie{
				Name:     cookie,
				Value:    "",
				Domain:   r.Host,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				// This is what deletes the cookie.
				MaxAge: -1,
			},
		)
	}

	// TODO: kill the refresh token.
}

func (a *API) GetLoggedInUser(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, user)
}
