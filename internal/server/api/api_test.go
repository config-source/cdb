package api

import (
	"net/http"

	"github.com/config-source/cdb/internal/auth"
	"github.com/config-source/cdb/internal/repository"
	"github.com/config-source/cdb/internal/services"
	"github.com/rs/zerolog"
)

func testAPI(
	repo repository.ModelRepository,
) (*API, *http.ServeMux, *auth.TestGateway) {
	gateway := auth.NewTestGateway()
	api, mux := New(
		repo,
		zerolog.New(nil).Level(zerolog.Disabled),
		[]byte("testing"),
		auth.NewUserService(
			gateway,
			gateway,
			true,
			"user-testing",
		),
		services.NewConfigValuesService(repo, true),
	)
	return api, mux, gateway
}
