package server

import (
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb/pkg/auth"
	"github.com/config-source/cdb/pkg/configkeys"
	"github.com/config-source/cdb/pkg/configvalues"
	"github.com/config-source/cdb/pkg/environments"
	"github.com/config-source/cdb/pkg/postgresutils"
	"github.com/config-source/cdb/pkg/services"
	"github.com/rs/zerolog"
)

func TestHealthCheckSuccess(t *testing.T) {
	gateway := auth.NewTestGateway()

	pool := postgresutils.InitTestDB(t)
	repoLogger := zerolog.New(nil).Level(zerolog.Disabled)

	svcRepo := services.NewRepository(repoLogger, pool)
	envRepo := environments.NewRepository(repoLogger, pool)
	keyRepo := configkeys.NewRepository(repoLogger, pool)
	valueRepo := configvalues.NewRepository(repoLogger, pool, envRepo)

	server := New(
		zerolog.New(nil).Level(zerolog.Disabled),
		[]byte("test key"),
		pool,
		auth.NewTestServiceWithGateway(gateway),
		configvalues.NewService(valueRepo, envRepo, keyRepo, gateway, true),
		environments.NewService(envRepo, gateway),
		configkeys.NewService(keyRepo, gateway),
		services.NewServiceService(svcRepo, gateway),
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
	gateway := auth.NewTestGateway()
	gateway.IsHealthy = false
	userService := auth.NewTestServiceWithGateway(gateway)

	pool := postgresutils.InitTestDB(t)
	repoLogger := zerolog.New(nil).Level(zerolog.Disabled)

	svcRepo := services.NewRepository(repoLogger, pool)
	envRepo := environments.NewRepository(repoLogger, pool)
	keyRepo := configkeys.NewRepository(repoLogger, pool)
	valueRepo := configvalues.NewRepository(repoLogger, pool, envRepo)

	server := New(
		zerolog.New(nil).Level(zerolog.Disabled),
		[]byte("test key"),
		pool,
		userService,
		configvalues.NewService(valueRepo, envRepo, keyRepo, gateway, true),
		environments.NewService(envRepo, gateway),
		configkeys.NewService(keyRepo, gateway),
		services.NewServiceService(svcRepo, gateway),
		"/frontend",
	)

	pool.Close()

	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	if rr.Code == 200 {
		t.Fatalf("Expected non-200 status code got: %d", rr.Code)
	}
}
