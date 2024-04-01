package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/config-source/cdb"
	"github.com/rs/zerolog"
)

var (
	ErrNotFound = errors.New("not found")
)

type Server struct {
	repo ModelRepository
	mux  *http.ServeMux
	log  zerolog.Logger
}

type ModelRepository interface {
	cdb.EnvironmentRepository
	cdb.ConfigValueRepository
	cdb.ConfigKeyRepository
	Healthy(context.Context) bool
}

func New(repo ModelRepository, log zerolog.Logger) *Server {
	server := &Server{
		repo: repo,
		log:  log,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/environments/by-name/{name}", server.GetEnvironmentByName)
	mux.HandleFunc("GET /healthz", server.HealtCheck)

	server.mux = mux
	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (s *Server) sendJson(w http.ResponseWriter, payload interface{}) {
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		s.log.Err(err).Msg("failed to encode a payload")
	}
}

func (s *Server) errorResponse(w http.ResponseWriter, message error) {
	response := ErrorResponse{
		Message: message.Error(),
	}
	s.sendJson(w, response)
}
