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

	cv, err := a.repo.GetConfigurationValue(r.Context(), environmentName, configKey)
	if err != nil {
		a.errorResponse(w, err)
		return
	}

	a.sendJson(w, cv)
}

func (a *API) SetConfigurationValue(w http.ResponseWriter, r *http.Request) {
	environmentName := r.PathValue("environment")
	configKey := r.PathValue("key")
	if environmentName == "" || configKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		a.errorResponse(w, errors.New("environment or key were empty"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var newConfigValue cdb.ConfigValue
	err := decoder.Decode(&newConfigValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.errorResponse(w, err)
		return
	}

	cv, err := a.configValueService.SetConfigurationValue(
		r.Context(),
		environmentName,
		configKey,
		newConfigValue,
	)
	if err != nil {
		a.errorResponse(w, err)
		return
	}

	a.sendJson(w, cv)
}

func (a *API) CreateConfigValue(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var newConfigValue cdb.ConfigValue
	err := decoder.Decode(&newConfigValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.errorResponse(w, err)
		return
	}

	newConfigValue, err = a.configValueService.CreateConfigValue(r.Context(), newConfigValue)
	if err != nil {
		a.errorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	a.sendJson(w, newConfigValue)
}

func (a *API) GetConfiguration(w http.ResponseWriter, r *http.Request) {
	environmentName := r.PathValue("environment")
	if environmentName == "" {
		w.WriteHeader(http.StatusBadRequest)
		a.errorResponse(w, errors.New("environment must not be blank"))
		return
	}

	cv, err := a.repo.GetConfiguration(r.Context(), environmentName)
	if err != nil {
		a.errorResponse(w, err)
		return
	}

	a.sendJson(w, cv)
}
