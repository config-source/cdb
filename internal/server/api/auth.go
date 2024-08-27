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
	// TODO: should be token set here
	token, err := auth.GenerateIdToken(a.tokenSigningKey, user)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	authCookies := []*http.Cookie{
		&http.Cookie{
			Name:     middleware.SessionCookieName,
			Value:    token,
			Domain:   r.Host,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		},
	}

	for _, cookie := range authCookies {
		http.SetCookie(w, cookie)
	}

	a.sendJson(w, auth.TokenSet{
		IDToken: token,
	})
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
