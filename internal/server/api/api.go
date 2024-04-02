package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/config-source/cdb/internal/repository"
	"github.com/rs/zerolog"
)

var (
	ErrNotFound = errors.New("not found")
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

func (a *API) errorResponse(w http.ResponseWriter, message error) {
	response := ErrorResponse{
		Message: message.Error(),
	}

	a.sendJson(w, response)
}
