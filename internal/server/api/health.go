package api

import (
	"net/http"
)

func (s *API) HealtCheck(w http.ResponseWriter, r *http.Request) {
	if !s.repo.Healthy(r.Context()) {
		w.WriteHeader(400)
	}

	w.Write([]byte{}) // nolint:errcheck
}
