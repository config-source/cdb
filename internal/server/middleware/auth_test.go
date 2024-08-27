package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/config-source/cdb/internal/auth"
	"github.com/config-source/cdb/internal/server/middleware"
	"github.com/rs/zerolog"
)

var signingKey []byte = []byte("testing")

func TestAuthenticationMiddlewareNoToken(t *testing.T) {
	var capturedUser *auth.User = nil
	testHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			capturedUser = middleware.GetUser(r)
		},
	)
	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	handler := middleware.Authentication(zerolog.New(os.Stdout), testHandler, signingKey)
	handler.ServeHTTP(rr, req)

	if capturedUser != nil {
		t.Error("Expected no user to be available due to lack of auth information!")
	}
}

func TestAuthenticationMiddlewareTokenInHeader(t *testing.T) {
	var capturedUser *auth.User = nil
	testHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			capturedUser = middleware.GetUser(r)
		},
	)
	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	expectedUser := auth.User{
		ID:    1,
		Email: "test@example.com",
	}
	token, err := auth.GenerateIdToken(signingKey, expectedUser)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	handler := middleware.Authentication(zerolog.New(os.Stdout), testHandler, signingKey)
	handler.ServeHTTP(rr, req)

	if !reflect.DeepEqual(*capturedUser, expectedUser) {
		t.Errorf("Expected user %v got %v", expectedUser, capturedUser)
	}
}

func TestAuthenticationMiddlewareTokenInvalidHeader(t *testing.T) {
	testHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
		},
	)
	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	expectedUser := auth.User{
		ID:    1,
		Email: "test@example.com",
	}
	token, err := auth.GenerateIdToken(signingKey, expectedUser)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("JWT %s", token))

	handler := middleware.Authentication(zerolog.New(os.Stdout), testHandler, signingKey)
	handler.ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d got %d", http.StatusBadRequest, rr.Result().StatusCode)
	}
}

func TestAuthenticationMiddlewareTokenInvalidToken(t *testing.T) {
	testHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
		},
	)
	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	expectedUser := auth.User{
		ID:    1,
		Email: "test@example.com",
	}
	token, err := auth.GenerateIdToken(signingKey, expectedUser)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token[10:]))

	handler := middleware.Authentication(zerolog.New(os.Stdout), testHandler, signingKey)
	handler.ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d got %d", http.StatusBadRequest, rr.Result().StatusCode)
	}
}

func TestAuthenticationMiddlewareTokenInCookie(t *testing.T) {
	var capturedUser *auth.User = nil
	testHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			capturedUser = middleware.GetUser(r)
		},
	)
	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	expectedUser := auth.User{
		ID:    1,
		Email: "test@example.com",
	}
	token, err := auth.GenerateIdToken(signingKey, expectedUser)
	if err != nil {
		t.Fatal(err)
	}

	req.AddCookie(&http.Cookie{
		Name:  middleware.SessionCookieName,
		Value: token,
	})

	handler := middleware.Authentication(zerolog.New(os.Stdout), testHandler, signingKey)
	handler.ServeHTTP(rr, req)

	if !reflect.DeepEqual(*capturedUser, expectedUser) {
		t.Errorf("Expected user %v got %v", expectedUser, capturedUser)
	}
}

func TestAuthenticationRequiredMiddlewareAllowsHandlerWithValidToken(t *testing.T) {
	called := false
	testHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			called = true
		},
	)
	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	user := auth.User{
		ID:    1,
		Email: "test@example.com",
	}
	token, err := auth.GenerateIdToken(signingKey, user)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	handler := middleware.AuthenticationRequired(zerolog.New(os.Stdout), testHandler, signingKey)
	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("Expected handler to be called because it had valid auth info!")
	}
}

func TestAuthenticationRequiredMiddlewareProtectsHandlerWithoutToken(t *testing.T) {
	called := false
	testHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			called = true
		},
	)
	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	handler := middleware.AuthenticationRequired(zerolog.New(os.Stdout), testHandler, signingKey)
	handler.ServeHTTP(rr, req)

	if called {
		t.Error("Expected handler to not be called because there was no valid auth info!")
	}
}
