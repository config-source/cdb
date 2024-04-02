package api

import (
	"errors"
	"net/http"

	"github.com/config-source/cdb"
)

func (a *API) GetEnvironmentByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		a.errorResponse(w, errors.New("name missing from url"))
		return
	}

	env, err := a.repo.GetEnvironmentByName(r.Context(), name)
	if err != nil {
		switch err {
		case cdb.ErrEnvNotFound:
			w.WriteHeader(http.StatusNotFound)
			err = ErrNotFound
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		a.errorResponse(w, err)
		return
	}

	a.sendJson(w, env)
}
