package api

import (
	"encoding/json"
	"net/http"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/internal/repository"
	"github.com/rs/zerolog"
)

type API struct {
	repo repository.ModelRepository
	log  zerolog.Logger
}

func New(repo repository.ModelRepository, log zerolog.Logger, mux *http.ServeMux) *API {
	api := &API{
		repo: repo,
		log:  log,
	}

	mux.HandleFunc("GET /api/v1/environments/by-name/{name}", api.GetEnvironmentByName)
	mux.HandleFunc("GET /api/v1/environments/by-id/{id}", api.GetEnvironmentByID)
	mux.HandleFunc("GET /api/v1/environments/tree", api.GetEnvironmentTree)
	mux.HandleFunc("POST /api/v1/environments", api.CreateEnvironment)

	mux.HandleFunc("POST /api/v1/config-keys", api.CreateConfigKey)
	mux.HandleFunc("GET /api/v1/config-keys", api.ListConfigKeys)
	mux.HandleFunc("GET /api/v1/config-keys/{id}", api.GetConfigKeyByID)

	mux.HandleFunc("POST /api/v1/config-values", api.CreateConfigValue)
	mux.HandleFunc("GET /api/v1/config-values/{environment}/{key}", api.GetConfigurationValue)
	mux.HandleFunc("GET /api/v1/config-values/{environment}", api.GetConfiguration)

	mux.HandleFunc("GET /healthz", api.HealtCheck)

	return api
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (a *API) sendJson(w http.ResponseWriter, payload interface{}) {
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		a.log.Err(err).Msg("failed to encode a payload")
	}
}

func (a *API) errorResponse(w http.ResponseWriter, err error) {
	switch err {
	case cdb.ErrEnvNotFound, cdb.ErrConfigKeyNotFound, cdb.ErrConfigValueNotFound:
		w.WriteHeader(http.StatusNotFound)
	// This is safe because subsequent calls to WriteHeader are ignored so
	// callers can set the status code before calling errorResponse but if they
	// haven't we want to send a 500.
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	response := ErrorResponse{
		Message: err.Error(),
	}

	a.sendJson(w, response)
}
