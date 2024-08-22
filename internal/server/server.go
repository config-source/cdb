package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/config-source/cdb/internal/auth"
	"github.com/config-source/cdb/internal/configvalues"
	"github.com/config-source/cdb/internal/repository"
	"github.com/config-source/cdb/internal/server/api"
	"github.com/rs/zerolog"
)

type Server struct {
	mux *http.ServeMux
	api *api.API
	ui  http.Handler
}

func New(
	repo repository.ModelRepository,
	userService *auth.UserService,
	configValueService *configvalues.Service,
	log zerolog.Logger,
	frontendLocation string,
) *Server {
	mux := http.NewServeMux()
	apiServer := api.New(repo, userService, configValueService, log, mux)

	var frontendHandler http.Handler
	if upstream, err := url.Parse(frontendLocation); err == nil && upstream.Scheme != "" {
		log.Info().Str("upstream", upstream.String()).Msg("Serving frontend from")
		frontendHandler = httputil.NewSingleHostReverseProxy(upstream)
	} else {
		log.Info().Str("location", frontendLocation).Msg("Serving frontend from")
		frontendHandler = http.FileServer(http.Dir(frontendLocation))
	}
	mux.Handle("GET /", frontendHandler)

	return &Server{
		mux: mux,
		api: apiServer,
		ui:  frontendHandler,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
