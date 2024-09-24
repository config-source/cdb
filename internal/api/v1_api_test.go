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
	"github.com/config-source/cdb/pkg/services"
	"github.com/config-source/cdb/pkg/testutils"
	"github.com/rs/zerolog"
)

func testAPI(
	repo *testutils.TestRepository,
	alwaysAuthd bool,
) (*V1, http.Handler, *auth.TestGateway) {
	gateway := auth.NewTestGateway()
	tokenSigningKey := []byte("test key")

	api, mux := NewV1(
		zerolog.New(nil).Level(zerolog.Disabled),
		tokenSigningKey,
		auth.NewTestServiceWithGateway(gateway),
		configvalues.NewService(repo, repo, repo, gateway, true),
		environments.NewService(repo, gateway),
		configkeys.NewService(repo, gateway),
		services.NewServiceService(repo, gateway),
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
	_, mux, _ := testAPI(&testutils.TestRepository{}, false)
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
