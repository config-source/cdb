package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb/internal/middleware"
	"github.com/config-source/cdb/pkg/auth"
)

func getCookieValue(cookies []*http.Cookie, cookieName string) string {
	for _, cookie := range cookies {
		if cookie.Name == cookieName {
			return cookie.Value
		}
	}

	return ""
}

func TestLogin(t *testing.T) {
	tc, mux := testAPI(t, true)
	gateway := tc.gateway

	_, err := gateway.CreateUser(context.Background(), auth.User{
		Email:    "test@example.com",
		Password: "Testing123!@",
	})
	if err != nil {
		t.Error(err)
	}

	creds := Credentials{
		Email:    "test@example.com",
		Password: "Testing123!@",
	}

	marshalled, err := json.Marshal(creds)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(marshalled))
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}

	var tokens auth.TokenSet
	if err := json.NewDecoder(rr.Body).Decode(&tokens); err != nil {
		t.Error(err)
	}

	if tokens.IDToken == "" {
		t.Errorf("Expected an ID token to be set!")
	}

	if tokens.AccessToken == "" {
		t.Errorf("Expected an Access token to be set!")
	}

	if tokens.RefreshToken == "" {
		t.Errorf("Expected an Refresh token to be set!")
	}

	cookies := rr.Result().Cookies()
	idTokenCookie := getCookieValue(cookies, middleware.IDTokenCookieName)
	accessTokenCookie := getCookieValue(cookies, middleware.AccessTokenCookieName)
	refreshTokenCookie := getCookieValue(cookies, middleware.RefreshTokenCookieName)

	if tokens.IDToken != idTokenCookie {
		t.Errorf("Expected tokenset IDToken and IDTokenCookie to match: %s %s", tokens.IDToken, idTokenCookie)
	}

	if tokens.AccessToken != accessTokenCookie {
		t.Errorf("Expected tokenset AccessToken and AccessTokenCookie to match: %s %s", tokens.AccessToken, accessTokenCookie)
	}

	if tokens.RefreshToken != refreshTokenCookie {
		t.Errorf("Expected tokenset RefreshToken and RefreshTokenCookie to match: %s %s", tokens.RefreshToken, refreshTokenCookie)
	}
}

func TestRegister(t *testing.T) {
	tc, mux := testAPI(t, true)
	gateway := tc.gateway

	creds := Credentials{
		Email:    "test@example.com",
		Password: "Testing123!@",
	}

	marshalled, err := json.Marshal(creds)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(marshalled))
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected status code 200 got: %d %s", rr.Code, rr.Body.String())
	}

	user, err := gateway.GetUser(context.Background(), auth.UserID(1))
	if err != nil {
		t.Error(err)
	}

	if user.Email != creds.Email {
		t.Errorf("Expected user email and creds email to match: %s %s", user.Email, creds.Email)
	}

	if user.Password != creds.Password {
		t.Errorf("Expected user password and creds password to match: %s %s", user.Password, creds.Password)
	}

	var tokens auth.TokenSet
	if err := json.NewDecoder(rr.Body).Decode(&tokens); err != nil {
		t.Error(err)
	}

	if tokens.IDToken == "" {
		t.Errorf("Expected an ID token to be set!")
	}

	if tokens.AccessToken == "" {
		t.Errorf("Expected an Access token to be set!")
	}

	if tokens.RefreshToken == "" {
		t.Errorf("Expected an Refresh token to be set!")
	}

	cookies := rr.Result().Cookies()
	idTokenCookie := getCookieValue(cookies, middleware.IDTokenCookieName)
	accessTokenCookie := getCookieValue(cookies, middleware.AccessTokenCookieName)
	refreshTokenCookie := getCookieValue(cookies, middleware.RefreshTokenCookieName)

	if tokens.IDToken != idTokenCookie {
		t.Errorf("Expected tokenset IDToken and IDTokenCookie to match: %s %s", tokens.IDToken, idTokenCookie)
	}

	if tokens.AccessToken != accessTokenCookie {
		t.Errorf("Expected tokenset AccessToken and AccessTokenCookie to match: %s %s", tokens.AccessToken, accessTokenCookie)
	}

	if tokens.RefreshToken != refreshTokenCookie {
		t.Errorf("Expected tokenset RefreshToken and RefreshTokenCookie to match: %s %s", tokens.RefreshToken, refreshTokenCookie)
	}
}

func TestRegisterWhenPublicRegistrationDisabled(t *testing.T) {
	tc, mux := testAPI(t, true)
	tc.api.userService = auth.NewUserService(
		tc.gateway,
		tc.gateway,
		&auth.TokenRegistry{},
		false,
		"",
	)

	creds := Credentials{
		Email:    "test@example.com",
		Password: "Testing123!@",
	}

	marshalled, err := json.Marshal(creds)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(marshalled))
	rr := httptest.NewRecorder()
	rr.Body = bytes.NewBuffer([]byte{})

	mux.ServeHTTP(rr, req)

	if rr.Code != 400 {
		t.Fatalf("Expected status code 400 got: %d %s", rr.Code, rr.Body.String())
	}
}
