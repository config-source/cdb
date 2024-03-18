package api

import (
	"errors"
	"net/http"
)

func (s *Server) GetEnvironmentByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse(w, errors.New("name missing from url"))
		return
	}

	env, err := s.repo.GetEnvironmentByName(r.Context(), name)
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			w.WriteHeader(http.StatusNotFound)
			err = ErrNotFound
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		errorResponse(w, err)
		return
	}

	sendJson(w, env)
}
