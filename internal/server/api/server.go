package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/config-source/cdb"
)

var (
	ErrNotFound = errors.New("not found")
)

type Server struct {
	repo ModelRepository
	mux  *http.ServeMux
}

type ModelRepository interface {
	cdb.EnvironmentRepository
	cdb.ConfigValueRepository
	cdb.ConfigKeyRepository
}

func New(repo ModelRepository) *Server {
	server := &Server{
		repo: repo,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/environments/by-name/{name}", server.GetEnvironmentByName)

	server.mux = mux
	return server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	s.mux.ServeHTTP(w, r)
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
