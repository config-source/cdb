package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/config-source/cdb/internal/middleware"
	"github.com/config-source/cdb/pkg/auth"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *V1) doLogin(w http.ResponseWriter, r *http.Request, user auth.User) {
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
				Path:     "/",
				SameSite: http.SameSiteStrictMode,
			},
		)
	}

	a.sendJson(w, tokens)
}

func (a *V1) Login(w http.ResponseWriter, r *http.Request) {
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

func (a *V1) Register(w http.ResponseWriter, r *http.Request) {
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

func (a *V1) Logout(w http.ResponseWriter, r *http.Request) {
	authCookies := []string{
		middleware.IDTokenCookieName,
		middleware.AccessTokenCookieName,
		middleware.RefreshTokenCookieName,
	}

	refreshToken := r.Header.Get("X-Refresh-Token")
	if refreshToken == "" {
		cookie, err := r.Cookie(middleware.RefreshTokenCookieName)
		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			a.sendErr(w, err)
			return
		} else if err == nil {
			refreshToken = cookie.Value
		}
	}

	if refreshToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, errors.New("no refresh token supplied"))
		return
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

	err := a.userService.Logout(r.Context(), refreshToken)
	if err != nil {
		a.sendErr(w, err)
	}
}

func (a *V1) GetLoggedInUser(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, user)
}

func (a *V1) IssueAPIToken(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	token, err := a.userService.IssueAPIToken(r.Context(), a.tokenSigningKey, user)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, token)
}

func (a *V1) ListAPITokens(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	tokens, err := a.userService.ListAPITokens(r.Context(), user)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, tokens)
}
