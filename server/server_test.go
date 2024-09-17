package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/auth"
	"github.com/config-source/cdb/configkeys"
	"github.com/config-source/cdb/configvalues"
	"github.com/config-source/cdb/environments"
	"github.com/config-source/cdb/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func testServer(repo *cdb.TestRepository) *http.ServeMux {
	gateway := auth.NewTestGateway()

	server := New(
		zerolog.New(nil).Level(zerolog.Disabled),
		[]byte("test key"),
		nil,
		auth.NewTestServiceWithGateway(gateway),
		configvalues.NewService(repo, repo, repo, gateway, true),
		environments.NewService(repo, gateway),
		configkeys.NewService(repo, gateway),
		services.NewServiceService(repo, gateway),
		"/frontend",
	)

	return server.mux
}

func TestHealthCheckSuccess(t *testing.T) {
	gateway := auth.NewTestGateway()

	pool, err := pgxpool.New(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}

	repo := &cdb.TestRepository{}
	server := New(
		zerolog.New(nil).Level(zerolog.Disabled),
		[]byte("test key"),
		pool,
		auth.NewTestServiceWithGateway(gateway),
		configvalues.NewService(repo, repo, repo, gateway, true),
		environments.NewService(repo, gateway),
		configkeys.NewService(repo, gateway),
		services.NewServiceService(repo, gateway),
		"/frontend",
	)

	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	server.mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected 200 status code got: %d", rr.Code)
	}
}

func TestHealthCheckFailure(t *testing.T) {
	repo := &cdb.TestRepository{
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
