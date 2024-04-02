package server

import (
	"net/http"

	"github.com/config-source/cdb/internal"
	"github.com/config-source/cdb/internal/server/api"
	"github.com/rs/zerolog"
)

type Server struct {
	mux *http.ServeMux
	api *api.API
}

func New(repo internal.ModelRepository, log zerolog.Logger) *Server {
	mux := http.NewServeMux()
	apiServer := api.New(repo, log, mux)

	return &Server{
		mux: mux,
		api: apiServer,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
