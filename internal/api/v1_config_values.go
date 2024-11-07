package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/config-source/cdb/internal/middleware"
	"github.com/config-source/cdb/pkg/configvalues"
)

func (a *V1) GetConfigurationValue(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	environmentID, err := strconv.Atoi(r.PathValue("environment"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, err)
		return
	}

	configKey := r.PathValue("key")
	if configKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, errors.New("key was empty"))
		return
	}

	cv, err := a.configValueService.GetConfigurationValue(r.Context(), user, environmentID, configKey)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	a.sendJson(w, cv)
}

func (a *V1) SetConfigurationValue(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	environmentID, err := strconv.Atoi(r.PathValue("environment"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, err)
		return
	}

	configKey := r.PathValue("key")
	if configKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, errors.New("key was empty"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var newConfigValue *configvalues.ConfigValue
	err = decoder.Decode(&newConfigValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, err)
		return
	}

	cv, err := a.configValueService.SetConfigurationValue(
		r.Context(),
		user,
		environmentID,
		configKey,
		newConfigValue,
	)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	a.sendJson(w, cv)
}

func (a *V1) CreateConfigValue(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var newConfigValue configvalues.ConfigValue
	err = decoder.Decode(&newConfigValue)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, err)
		return
	}

	newConfigValue, err = a.configValueService.CreateConfigValue(r.Context(), user, newConfigValue)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	a.sendJson(w, newConfigValue)
}

func (a *V1) GetConfiguration(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	environmentID, err := strconv.Atoi(r.PathValue("environment"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, err)
		return
	}

	cv, err := a.configValueService.GetConfiguration(r.Context(), user, environmentID)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	a.sendJson(w, cv)
}
