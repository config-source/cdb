package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/auth"
	"github.com/config-source/cdb/configkeys"
	"github.com/config-source/cdb/configvalues"
	"github.com/config-source/cdb/environments"
	"github.com/config-source/cdb/server/middleware"
	"github.com/config-source/cdb/services"
	"github.com/rs/zerolog"
)

type API struct {
	log             zerolog.Logger
	tokenSigningKey []byte

	userService        *auth.UserService
	configValueService *configvalues.Service
	envService         *environments.Service
	svcService         *services.ServiceService
	configKeyService   *configkeys.Service
}

func New(
	log zerolog.Logger,
	tokenSigningKey []byte,
	userService *auth.UserService,
	configValueService *configvalues.Service,
	envService *environments.Service,
	configKeyService *configkeys.Service,
	svcService *services.ServiceService,
) (*API, http.Handler) {
	api := &API{
		log:             log,
		tokenSigningKey: tokenSigningKey,

		configValueService: configValueService,
		envService:         envService,
		configKeyService:   configKeyService,
		userService:        userService,
		svcService:         svcService,
	}

	// v1 routes
	v1Mux := http.NewServeMux()

	v1Mux.HandleFunc("GET /api/v1/environments/by-name/{name}", api.GetEnvironmentByName)
	v1Mux.HandleFunc("GET /api/v1/environments/by-id/{id}", api.GetEnvironmentByID)
	v1Mux.HandleFunc("GET /api/v1/environments/tree", api.GetEnvironmentTree)
	v1Mux.HandleFunc("GET /api/v1/environments", api.ListEnvironments)
	v1Mux.HandleFunc("POST /api/v1/environments", api.CreateEnvironment)

	v1Mux.HandleFunc("GET /api/v1/services/by-name/{name}", api.GetServiceByName)
	v1Mux.HandleFunc("GET /api/v1/services/by-id/{id}", api.GetServiceByID)
	v1Mux.HandleFunc("GET /api/v1/services", api.ListServices)
	v1Mux.HandleFunc("POST /api/v1/services", api.CreateService)

	v1Mux.HandleFunc("POST /api/v1/config-keys", api.CreateConfigKey)
	v1Mux.HandleFunc("GET /api/v1/config-keys", api.ListConfigKeys)
	v1Mux.HandleFunc("GET /api/v1/config-keys/by-id/{id}", api.GetConfigKeyByID)
	v1Mux.HandleFunc("GET /api/v1/config-keys/by-name/{name}", api.GetConfigKeyByName)

	v1Mux.HandleFunc("POST /api/v1/config-values", api.CreateConfigValue)
	v1Mux.HandleFunc("GET /api/v1/config-values/{environment}/{key}", api.GetConfigurationValue)
	v1Mux.HandleFunc("POST /api/v1/config-values/{environment}/{key}", api.SetConfigurationValue)
	v1Mux.HandleFunc("GET /api/v1/config-values/{environment}", api.GetConfiguration)

	v1Mux.HandleFunc("GET /api/v1/users/me", api.GetLoggedInUser)
	v1Mux.HandleFunc("POST /api/v1/auth/api-tokens", api.IssueAPIToken)
	v1Mux.HandleFunc("GET /api/v1/auth/api-tokens", api.ListAPITokens)

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("DELETE /api/v1/auth/logout", api.Logout)
	apiMux.HandleFunc("POST /api/v1/auth/login", api.Login)
	apiMux.HandleFunc("POST /api/v1/auth/register", api.Register)
	apiMux.Handle("/api/v1/", middleware.AuthenticationRequired(log, userService, tokenSigningKey, v1Mux))

	return api, apiMux
}

func (a *API) sendJson(w http.ResponseWriter, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		a.log.Err(err).Msg("failed to encode a payload")
	}
}

func (a *API) sendErr(w http.ResponseWriter, err error) {
	switch {
	case
		errors.Is(err, auth.ErrUserNotFound),
		errors.Is(err, environments.ErrNotFound),
		errors.Is(err, configkeys.ErrNotFound),
		errors.Is(err, configvalues.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case
		errors.Is(err, configvalues.ErrNotValid),
		errors.Is(err, configvalues.ErrAlreadySet),
		errors.Is(err, auth.ErrPublicRegisterDisabled),
		errors.Is(err, auth.ErrEmailInUse):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, auth.ErrUnauthorized),
		errors.Is(err, auth.ErrInvalidPassword):
		w.WriteHeader(http.StatusForbidden)
	case errors.Is(err, auth.ErrUnauthenticated):
		w.WriteHeader(http.StatusUnauthorized)
	// This is safe because subsequent calls to WriteHeader are ignored so
	// callers can set the status code before calling errorResponse but if they
	// haven't we want to send a 500.
	default:
		a.log.Err(err).Msg("unhandled error")
		w.WriteHeader(http.StatusInternalServerError)
	}

	a.sendJson(w, cdb.NewErrorResponse(err.Error()))
}
