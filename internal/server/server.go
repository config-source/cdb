package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/config-source/cdb/internal/auth"
	"github.com/config-source/cdb/internal/repository"
	"github.com/config-source/cdb/internal/server/api"
	"github.com/config-source/cdb/internal/services"
	"github.com/rs/zerolog"
)

type Server struct {
	mux *http.ServeMux
	api *api.API
	ui  http.Handler
}

func New(
	repo repository.ModelRepository,
	log zerolog.Logger,
	tokenSigningKey []byte,
	userService *auth.UserService,
	configValueService *services.ConfigValuesService,
	frontendLocation string,
) *Server {
	var frontendHandler http.Handler
	frontendServingLog := log.Info()
	defer frontendServingLog.Msg("Serving frontend from")

	if upstream, err := url.Parse(frontendLocation); err == nil && upstream.Scheme != "" {
		frontendServingLog.Str("upstream", upstream.String())
		frontendHandler = httputil.NewSingleHostReverseProxy(upstream)
	} else {
		frontendServingLog.Str("location", frontendLocation)
		frontendHandler = http.FileServer(http.Dir(frontendLocation))
	}

	apiServer, apiMux := api.New(
		repo,
		log,
		tokenSigningKey,
		userService,
		configValueService,
	)

	apiMux.Handle("/", frontendHandler)

	return &Server{
		mux: apiMux,
		api: apiServer,
		ui:  frontendHandler,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
