package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/config-source/cdb/internal/middleware"
	"github.com/config-source/cdb/pkg/configkeys"
)

func (a *V1) GetConfigKeyByID(w http.ResponseWriter, r *http.Request) {
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

	ck, err := a.configKeyService.GetConfigKeyByID(r.Context(), user, id)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, ck)
}

func (a *V1) GetConfigKeyByName(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	name := r.PathValue("name")
	serviceName := r.PathValue("serviceName")

	ck, err := a.configKeyService.GetConfigKeyByName(r.Context(), user, serviceName, name)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, ck)
}

func (a *V1) ListConfigKeys(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	serviceNames := r.URL.Query()["service"]

	cks, err := a.configKeyService.ListConfigKeys(r.Context(), user, serviceNames...)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, cks)
}

func (a *V1) CreateConfigKey(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var configKey configkeys.ConfigKey
	err = decoder.Decode(&configKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, err)
		return
	}

	configKey, err = a.configKeyService.CreateConfigKey(r.Context(), user, configKey)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	a.sendJson(w, configKey)
}
