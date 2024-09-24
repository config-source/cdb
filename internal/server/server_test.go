package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb/pkg/auth"
	"github.com/config-source/cdb/pkg/configkeys"
	"github.com/config-source/cdb/pkg/configvalues"
	"github.com/config-source/cdb/pkg/environments"
	"github.com/config-source/cdb/pkg/services"
	"github.com/config-source/cdb/pkg/testutils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func testServer(repo *testutils.TestRepository) http.Handler {
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

	return server.handler
}

func TestHealthCheckSuccess(t *testing.T) {
	gateway := auth.NewTestGateway()

	pool, err := pgxpool.New(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}

	repo := &testutils.TestRepository{}
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

	server.handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected 200 status code got: %d", rr.Code)
	}
}

func TestHealthCheckFailure(t *testing.T) {
	repo := &testutils.TestRepository{
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
