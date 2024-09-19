package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/config-source/cdb/internal/server/middleware"
	"github.com/config-source/cdb/pkg/configvalues"
)

func (a *API) GetConfigurationValue(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	environmentName := r.PathValue("environment")
	configKey := r.PathValue("key")
	if environmentName == "" || configKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, errors.New("environment or key were empty"))
		return
	}

	cv, err := a.configValueService.GetConfigurationValue(r.Context(), user, environmentName, configKey)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, cv)
}

func (a *API) SetConfigurationValue(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	environmentName := r.PathValue("environment")
	configKey := r.PathValue("key")
	if environmentName == "" || configKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, errors.New("environment or key were empty"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var newConfigValue *configvalues.ConfigValue
	err = decoder.Decode(&newConfigValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, err)
		return
	}

	cv, err := a.configValueService.SetConfigurationValue(
		r.Context(),
		user,
		environmentName,
		configKey,
		newConfigValue,
	)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, cv)
}

func (a *API) CreateConfigValue(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var newConfigValue configvalues.ConfigValue
	err = decoder.Decode(&newConfigValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, err)
		return
	}

	newConfigValue, err = a.configValueService.CreateConfigValue(r.Context(), user, newConfigValue)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	a.sendJson(w, newConfigValue)
}

func (a *API) GetConfiguration(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	environmentName := r.PathValue("environment")
	if environmentName == "" {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, errors.New("environment must not be blank"))
		return
	}

	cv, err := a.configValueService.GetConfiguration(r.Context(), user, environmentName)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, cv)
}
