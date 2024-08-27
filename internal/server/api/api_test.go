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
) (*API, *http.ServeMux) {
	mux := http.NewServeMux()
	return New(
		repo,
		zerolog.New(nil).Level(zerolog.Disabled),
		[]byte("testing"),
		auth.NewUserService(
			auth.NewTestGateway(),
			auth.NewTestGateway(),
			true,
			"user-testing",
		),
		configvalues.NewService(repo, true),
		mux,
	), mux
}
