package api

import (
	"net/http"
)

func (s *Server) HealtCheck(w http.ResponseWriter, r *http.Request) {
	if !s.repo.Healthy(r.Context()) {
		w.WriteHeader(400)
	}

	w.Write([]byte{})
}
