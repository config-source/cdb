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
	log zerolog.Logger,
	tokenSigningKey []byte,
	userService *auth.UserService,
	configValueService *configvalues.Service,
	frontendLocation string,
) *Server {
	mux := http.NewServeMux()
	apiServer := api.New(
		repo,
		log,
		tokenSigningKey,
		userService,
		configValueService,
		mux,
	)

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
