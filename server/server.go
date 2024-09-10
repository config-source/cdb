package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/config-source/cdb/auth"
	"github.com/config-source/cdb/repository"
	"github.com/config-source/cdb/server/api"
	"github.com/config-source/cdb/services"
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
	configValueService *services.ConfigValues,
	envService *services.Environments,
	configKeysService *services.ConfigKeys,
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
		log,
		tokenSigningKey,
		userService,
		configValueService,
		envService,
		configKeysService,
	)

	mux := http.NewServeMux()
	mux.Handle("/api/", apiMux)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		if !repo.Healthy(r.Context()) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		if !userService.Healthy(r.Context()) {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		w.Write([]byte{}) // nolint:errcheck
	})
	mux.Handle("/", frontendHandler)

	return &Server{
		mux: mux,
		api: apiServer,
		ui:  frontendHandler,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
