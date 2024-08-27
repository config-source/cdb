package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/internal/auth"
	"github.com/config-source/cdb/internal/configvalues"
	"github.com/config-source/cdb/internal/repository"
	"github.com/rs/zerolog"
)

type API struct {
	repo repository.ModelRepository
	log  zerolog.Logger

	tokenSigningKey []byte

	userService        *auth.UserService
	configValueService *configvalues.Service
}

func New(
	repo repository.ModelRepository,
	log zerolog.Logger,
	tokenSigningKey []byte,
	userService *auth.UserService,
	configValueService *configvalues.Service,
	mux *http.ServeMux,
) *API {
	api := &API{
		repo: repo,
		log:  log,

		tokenSigningKey: tokenSigningKey,

		configValueService: configValueService,
		userService:        userService,
	}

	mux.HandleFunc("GET /api/v1/environments/by-name/{name}", api.GetEnvironmentByName)
	mux.HandleFunc("GET /api/v1/environments/by-id/{id}", api.GetEnvironmentByID)
	mux.HandleFunc("GET /api/v1/environments/tree", api.GetEnvironmentTree)
	mux.HandleFunc("GET /api/v1/environments", api.ListEnvironments)
	mux.HandleFunc("POST /api/v1/environments", api.CreateEnvironment)

	mux.HandleFunc("POST /api/v1/config-keys", api.CreateConfigKey)
	mux.HandleFunc("GET /api/v1/config-keys", api.ListConfigKeys)
	mux.HandleFunc("GET /api/v1/config-keys/by-id/{id}", api.GetConfigKeyByID)
	mux.HandleFunc("GET /api/v1/config-keys/by-name/{name}", api.GetConfigKeyByName)

	mux.HandleFunc("POST /api/v1/config-values", api.CreateConfigValue)
	mux.HandleFunc("GET /api/v1/config-values/{environment}/{key}", api.GetConfigurationValue)
	mux.HandleFunc("POST /api/v1/config-values/{environment}/{key}", api.SetConfigurationValue)
	mux.HandleFunc("GET /api/v1/config-values/{environment}", api.GetConfiguration)

	mux.HandleFunc("GET /healthz", api.HealtCheck)

	return api
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
		errors.Is(err, cdb.ErrEnvNotFound),
		errors.Is(err, cdb.ErrConfigKeyNotFound),
		errors.Is(err, cdb.ErrConfigValueNotFound):
		w.WriteHeader(http.StatusNotFound)
	case errors.Is(err, cdb.ErrConfigValueNotValid),
		errors.Is(err, cdb.ErrConfigValueAlreadySet),
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

	response := cdb.ErrorResponse{
		Message: err.Error(),
	}

	a.sendJson(w, response)
}
