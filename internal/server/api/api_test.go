package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb/internal/auth"
	"github.com/config-source/cdb/internal/repository"
	"github.com/config-source/cdb/internal/server/middleware"
	"github.com/config-source/cdb/internal/services"
	"github.com/rs/zerolog"
)

func testAPI(
	repo repository.ModelRepository,
	alwaysAuthd bool,
) (*API, http.Handler, *auth.TestGateway) {
	gateway := auth.NewTestGateway()
	tokenSigningKey := []byte("test key")

	api, mux := New(
		zerolog.New(nil).Level(zerolog.Disabled),
		tokenSigningKey,
		auth.NewUserService(
			gateway,
			gateway,
			true,
			"user-testing",
		),
		services.NewConfigValuesService(repo, gateway, true),
		services.NewEnvironmentsService(repo, gateway),
		services.NewConfigKeysService(repo, gateway),
	)

	if alwaysAuthd {
		return api, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		}), gateway
	}

	return api, mux, gateway
}

func TestProtectedRoutesAreProtected(t *testing.T) {
	_, mux, _ := testAPI(&repository.TestRepository{}, false)
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
