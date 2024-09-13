package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/config-source/cdb/auth"
	"github.com/config-source/cdb/configkeys"
	"github.com/config-source/cdb/configvalues"
	"github.com/config-source/cdb/environments"
	"github.com/config-source/cdb/postgresutils"
	"github.com/config-source/cdb/server/api"
	"github.com/config-source/cdb/services"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Server struct {
	mux *http.ServeMux
	api *api.API
	ui  http.Handler
}

func New(
	log zerolog.Logger,
	tokenSigningKey []byte,
	postgresPool *pgxpool.Pool,
	userService *auth.UserService,
	configValueService *configvalues.Service,
	envService *environments.Service,
	configKeyService *configkeys.Service,
	svcService *services.ServiceService,
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
		configKeyService,
		svcService,
	)

	mux := http.NewServeMux()
	mux.Handle("/api/", apiMux)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		if !postgresutils.HealthCheck(r.Context(), postgresPool, log) {
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
