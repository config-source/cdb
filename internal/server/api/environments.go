package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
		a.errorResponse(w, err)
		return
	}

	a.sendJson(w, env)
}

func (a *API) GetEnvironmentByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.errorResponse(w, err)
		return
	}

	env, err := a.repo.GetEnvironment(r.Context(), id)
	if err != nil {
		a.errorResponse(w, err)
		return
	}

	a.sendJson(w, env)
}

func (a *API) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var env cdb.Environment
	err := decoder.Decode(&env)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.errorResponse(w, err)
		return
	}

	env, err = a.repo.CreateEnvironment(r.Context(), env)
	if err != nil {
		a.errorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	a.sendJson(w, env)
}
