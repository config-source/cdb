package api

import (
	"net/http"

	"github.com/config-source/cdb/internal/repository"
	"github.com/rs/zerolog"
)

func testAPI(repo repository.ModelRepository) (*API, *http.ServeMux) {
	mux := http.NewServeMux()
	return New(repo, zerolog.New(nil).Level(zerolog.Disabled), mux), mux
}
