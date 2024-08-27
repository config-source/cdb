package api

import (
	"net/http"
)

func (s *API) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if !s.repo.Healthy(r.Context()) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if !s.userService.Healthy(r.Context()) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Write([]byte{}) // nolint:errcheck
}
