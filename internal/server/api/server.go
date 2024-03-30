package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wr := &StatusRecorder{ResponseWriter: w}
	wr.Header().Set("Content-Type", "application/json")
	s.mux.ServeHTTP(wr, r)
	accessLog := s.log.Info()
	if wr.Status >= 400 {
		accessLog = s.log.Error()
	}

	accessLog.
		Str("url", r.URL.String()).
		Int("statusCode", wr.Status).
		Msg("request served")
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func sendJson(w http.ResponseWriter, payload interface{}) {
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		// TODO: logging
		fmt.Println(err)
	}
}

func errorResponse(w http.ResponseWriter, message error) {
	response := ErrorResponse{
		Message: message.Error(),
	}
	sendJson(w, response)
}
