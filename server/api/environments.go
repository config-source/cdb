package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/config-source/cdb"
	"github.com/config-source/cdb/server/middleware"
)

func (a *API) GetEnvironmentByName(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	name := r.PathValue("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, errors.New("name missing from url"))
		return
	}

	env, err := a.envService.GetEnvironmentByName(r.Context(), user, name)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, env)
}

func (a *API) GetEnvironmentByID(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, err)
		return
	}

	env, err := a.envService.GetEnvironmentByID(r.Context(), user, id)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, env)
}

func (a *API) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var env cdb.Environment
	err = decoder.Decode(&env)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, err)
		return
	}

	env, err = a.envService.CreateEnvironment(r.Context(), user, env)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	a.sendJson(w, env)
}

func (a *API) GetEnvironmentTree(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	trees, err := a.envService.GetEnvironmentTree(r.Context(), user)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, trees)
}

func (a *API) ListEnvironments(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	environs, err := a.envService.ListEnvironments(r.Context(), user)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, environs)
}
