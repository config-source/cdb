package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/config-source/cdb"
)

func (a *API) GetConfigurationValue(w http.ResponseWriter, r *http.Request) {
	environmentName := r.PathValue("environment")
	configKey := r.PathValue("key")
	if environmentName == "" || configKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		a.errorResponse(w, errors.New("environment or key were empty"))
		return
	}

	ck, err := a.repo.GetConfigurationValue(r.Context(), environmentName, configKey)
	if err != nil {
		a.errorResponse(w, err)
		return
	}

	a.sendJson(w, ck)
}

func (a *API) CreateConfigValue(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var env cdb.ConfigValue
	err := decoder.Decode(&env)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.errorResponse(w, err)
		return
	}

	env, err = a.repo.CreateConfigValue(r.Context(), env)
	if err != nil {
		a.errorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	a.sendJson(w, env)
}

func (a *API) GetConfiguration(w http.ResponseWriter, r *http.Request) {
	environmentName := r.PathValue("environment")
	if environmentName == "" {
		w.WriteHeader(http.StatusBadRequest)
		a.errorResponse(w, errors.New("environment must not be blank"))
		return
	}

	ck, err := a.repo.GetConfiguration(r.Context(), environmentName)
	if err != nil {
		a.errorResponse(w, err)
		return
	}

	a.sendJson(w, ck)
}
