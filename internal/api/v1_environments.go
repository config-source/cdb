package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/config-source/cdb/internal/middleware"
	"github.com/config-source/cdb/pkg/environments"
)

func (a *V1) GetEnvironmentByName(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	name := r.PathValue("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, errors.New("name missing from url"))
		return
	}

	serviceName := r.PathValue("serviceName")
	if serviceName == "" {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, errors.New("serviceName missing from url"))
		return
	}

	env, err := a.envService.GetEnvironmentByName(r.Context(), user, serviceName, name)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	a.sendJson(w, env)
}

func (a *V1) GetEnvironmentByID(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, err)
		return
	}

	env, err := a.envService.GetEnvironmentByID(r.Context(), user, id)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	a.sendJson(w, env)
}

func (a *V1) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var env environments.Environment
	err = decoder.Decode(&env)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, err)
		return
	}

	env, err = a.envService.CreateEnvironment(r.Context(), user, env)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	a.sendJson(w, env)
}

func (a *V1) GetEnvironmentTree(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	trees, err := a.envService.GetEnvironmentTree(r.Context(), user)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	a.sendJson(w, trees)
}

func (a *V1) ListEnvironments(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	environs, err := a.envService.ListEnvironments(r.Context(), user)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	a.sendJson(w, environs)
}

func (a *V1) UpdateEnvironment(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var env environments.Environment
	err = decoder.Decode(&env)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, err)
		return
	}

	if env.ID != id {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, errors.New("cannot change id of environment"))
		return
	}

	updated, err := a.envService.UpdateEnvironment(r.Context(), user, env)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	a.sendJson(w, updated)
}

func (a *V1) DeleteEnvironment(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, err)
		return
	}

	err = a.envService.DeleteEnvironment(r.Context(), user, id)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	w.Write([]byte{})
}
