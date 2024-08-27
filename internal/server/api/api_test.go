package api

import (
	"net/http"

	"github.com/config-source/cdb/internal/auth"
	"github.com/config-source/cdb/internal/configvalues"
	"github.com/config-source/cdb/internal/repository"
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
		configvalues.NewService(repo, true),
	)
	return api, mux, gateway
}
