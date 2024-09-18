package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/config-source/cdb/configkeys"
	"github.com/config-source/cdb/server/middleware"
)

func (a *API) GetConfigKeyByID(w http.ResponseWriter, r *http.Request) {
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

	serviceID, err := strconv.Atoi(r.PathValue("serviceID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, err)
		return
	}

	ck, err := a.configKeyService.GetConfigKeyByID(r.Context(), user, serviceID, id)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, ck)
}

func (a *API) GetConfigKeyByName(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	name := r.PathValue("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, errors.New("name must be provided"))
		return
	}

	serviceID, err := strconv.Atoi(r.PathValue("serviceID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, err)
		return
	}

	ck, err := a.configKeyService.GetConfigKeyByName(r.Context(), user, serviceID, name)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, ck)
}

func (a *API) ListConfigKeys(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	var serviceIDs []int
	if svcIDStrings, ok := r.URL.Query()["service"]; ok {
		serviceIDs = make([]int, len(svcIDStrings))
		for i, svcIDString := range svcIDStrings {
			serviceIDs[i], err = strconv.Atoi(svcIDString)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				a.sendErr(w, err)
				return
			}
		}
	}

	cks, err := a.configKeyService.ListConfigKeys(r.Context(), user, serviceIDs...)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, cks)
}

func (a *API) CreateConfigKey(w http.ResponseWriter, r *http.Request) {
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
