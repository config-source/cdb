package api

import (
	"fmt"
	"net/http"

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

			fmt.Println("Setting auth header to:", idToken)
			r.Header.Set(
				"Authorization",
				fmt.Sprintf("%s%s", middleware.AuthorizationHeaderPrefix, idToken),
			)

			mux.ServeHTTP(w, r)
		}), gateway
	}

	return api, mux, gateway
}
