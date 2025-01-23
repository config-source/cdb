package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb/internal/middleware"
	"github.com/config-source/cdb/pkg/auth"
	"github.com/config-source/cdb/pkg/configkeys"
	"github.com/config-source/cdb/pkg/configvalues"
	"github.com/config-source/cdb/pkg/environments"
	"github.com/config-source/cdb/pkg/postgresutils"
	"github.com/config-source/cdb/pkg/services"
	"github.com/rs/zerolog"
)

type TestContext struct {
	serviceRepo     *services.Repository
	environmentRepo *environments.Repository
	keyRepo         *configkeys.Repository
	valueRepo       *configvalues.Repository

	api *V1

	gateway *auth.TestGateway
}

func testAPI(
	t *testing.T,
	alwaysAuthd bool,
) (TestContext, http.Handler) {
	t.Helper()

	gateway := auth.NewTestGateway()
	tokenSigningKey := []byte("test key")

	pool := postgresutils.InitTestDB(t)
	repoLogger := zerolog.New(nil).Level(zerolog.Disabled)

	svcRepo := services.NewRepository(repoLogger, pool)
	envRepo := environments.NewRepository(repoLogger, pool)
	keyRepo := configkeys.NewRepository(repoLogger, pool)
	valueRepo := configvalues.NewRepository(repoLogger, pool, envRepo)

	api, mux := NewV1(
		zerolog.New(nil).Level(zerolog.Disabled),
		tokenSigningKey,
		auth.NewTestServiceWithGateway(gateway),
		configvalues.NewService(valueRepo, envRepo, keyRepo, gateway, true),
		environments.NewService(envRepo, gateway),
		configkeys.NewService(keyRepo, gateway),
		services.NewServiceService(svcRepo, gateway),
	)

	tc := TestContext{
		serviceRepo:     svcRepo,
		environmentRepo: envRepo,
		keyRepo:         keyRepo,
		valueRepo:       valueRepo,
		api:             api,
		gateway:         gateway,
	}

	if alwaysAuthd {
		return tc, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := auth.User{}
			idToken, err := auth.GenerateIdToken(tokenSigningKey, user)
			if err != nil {
				panic(err)
			}

			r.Header.Set(
				"Authorization",
				fmt.Sprintf("%s%s", middleware.AuthorizationHeaderPrefix, idToken),
			)

			mux.ServeHTTP(w, r)
		})
	}

	return tc, mux
}

func TestProtectedRoutesAreProtected(t *testing.T) {
	_, mux := testAPI(t, false)
	protectedRoutes := []struct {
		endpoint string
		method   string
	}{
		{endpoint: "/api/v1/environments/by-name/test", method: "GET"},
		{endpoint: "/api/v1/environments/by-id/1", method: "GET"},
		{endpoint: "/api/v1/environments/tree", method: "GET"},
		{endpoint: "/api/v1/environments", method: "GET"},
		{endpoint: "/api/v1/environments", method: "POST"},

		{endpoint: "/api/v1/config-keys", method: "POST"},
		{endpoint: "/api/v1/config-keys", method: "GET"},
		{endpoint: "/api/v1/config-keys/by-id/1", method: "GET"},
		{endpoint: "/api/v1/config-keys/by-name/test", method: "GET"},

		{endpoint: "/api/v1/config-values", method: "POST"},
		{endpoint: "/api/v1/config-values/test/testKey", method: "GET"},
		{endpoint: "/api/v1/config-values/test/testKey", method: "POST"},
		{endpoint: "/api/v1/config-values/test", method: "GET"},
	}

	for _, route := range protectedRoutes {
		req := httptest.NewRequest(route.method, route.endpoint, nil)
		rr := httptest.NewRecorder()
		rr.Body = bytes.NewBuffer([]byte{})
		mux.ServeHTTP(rr, req)
		if rr.Code != 401 {
			t.Errorf("Expected 401 got: %d", rr.Code)
		}
	}
}
