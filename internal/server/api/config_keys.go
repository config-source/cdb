package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/config-source/cdb"
)

func (a *API) GetConfigKeyByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.errorResponse(w, err)
		return
	}

	ck, err := a.repo.GetConfigKey(r.Context(), id)
	if err != nil {
		a.errorResponse(w, err)
		return
	}

	a.sendJson(w, ck)
}

func (a *API) GetConfigKeyByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		a.errorResponse(w, errors.New("name must be provided"))
		return
	}

	ck, err := a.repo.GetConfigKeyByName(r.Context(), name)
	if err != nil {
		a.errorResponse(w, err)
		return
	}

	a.sendJson(w, ck)
}

func (a *API) ListConfigKeys(w http.ResponseWriter, r *http.Request) {
	cks, err := a.repo.ListConfigKeys(r.Context())
	if err != nil {
		a.errorResponse(w, err)
		return
	}

	a.sendJson(w, cks)
}

func (a *API) CreateConfigKey(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var env cdb.ConfigKey
	err := decoder.Decode(&env)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.errorResponse(w, err)
		return
	}

	env, err = a.repo.CreateConfigKey(r.Context(), env)
	if err != nil {
		a.errorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	a.sendJson(w, env)
}
