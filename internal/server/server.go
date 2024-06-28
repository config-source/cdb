package server

import (
	"net/http"

	"github.com/config-source/cdb/internal/configvalues"
	"github.com/config-source/cdb/internal/repository"
	"github.com/config-source/cdb/internal/server/api"
	"github.com/config-source/cdb/internal/server/ui"
	"github.com/rs/zerolog"
)

type Server struct {
	mux *http.ServeMux
	api *api.API
	ui  *ui.UI
}

func New(
	repo repository.ModelRepository,
	configValueService *configvalues.Service,
	log zerolog.Logger,
	frontendLocation string,
) *Server {
	apiMux := http.NewServeMux()
	apiServer := api.New(repo, configValueService, log, apiMux)

	uiMux := http.NewServeMux()
	uiServer := ui.New(repo, configValueService, log, uiMux)

	mux := http.NewServeMux()
	mux.Handle("/api", apiMux)
	mux.Handle("/", uiMux)

	return &Server{
		mux: mux,
		api: apiServer,
		ui:  uiServer,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
