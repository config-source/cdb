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
		a.sendErr(w, r, err)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, err)
		return
	}

	ck, err := a.configKeyService.GetConfigKeyByID(r.Context(), user, id)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	a.sendJson(w, ck)
}

func (a *V1) GetConfigKeyByName(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	name := r.PathValue("name")
	serviceName := r.PathValue("serviceName")

	ck, err := a.configKeyService.GetConfigKeyByName(r.Context(), user, serviceName, name)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	a.sendJson(w, ck)
}

func (a *V1) ListConfigKeys(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	services := r.URL.Query()["service"]
	serviceIDs := make([]int, len(services))
	for i, svc := range services {
		var err error
		serviceIDs[i], err = strconv.Atoi(svc)
		if err != nil {
			a.sendErr(w, r, err)
			return
		}
	}

	cks, err := a.configKeyService.ListConfigKeys(r.Context(), user, serviceIDs...)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	a.sendJson(w, cks)
}

func (a *V1) CreateConfigKey(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var configKey configkeys.ConfigKey
	err = decoder.Decode(&configKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, r, err)
		return
	}

	configKey, err = a.configKeyService.CreateConfigKey(r.Context(), user, configKey)
	if err != nil {
		a.sendErr(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	a.sendJson(w, configKey)
}
