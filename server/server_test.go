package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb/auth"
	"github.com/config-source/cdb/configkeys"
	"github.com/config-source/cdb/configvalues"
	"github.com/config-source/cdb/repository"
	"github.com/rs/zerolog"
)

func testServer(repo repository.ModelRepository) *http.ServeMux {
	gateway := auth.NewTestGateway()

	server := New(
		repo,
		zerolog.New(nil).Level(zerolog.Disabled),
		[]byte("test key"),
		auth.NewUserService(
			gateway,
			gateway,
			true,
			"user-testing",
		),
		configvalues.NewService(repo, gateway, true),
		environments.NewService(repo, gateway),
		configkeys.NewService(repo, gateway),
		"/frontend",
	)

	return server.mux
}

func TestHealthCheckSuccess(t *testing.T) {
	repo := &repository.TestRepository{
		IsHealthy: true,
	}

	mux := testServer(repo)
	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected 200 status code got: %d", rr.Code)
	}
}

func TestHealthCheckFailure(t *testing.T) {
	repo := &repository.TestRepository{
		IsHealthy: false,
	}

	mux := testServer(repo)
	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code == 200 {
		t.Fatalf("Expected non-200 status code got: %d", rr.Code)
	}
}
