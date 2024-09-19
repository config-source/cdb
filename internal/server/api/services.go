package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/config-source/cdb/internal/server/middleware"
	"github.com/config-source/cdb/pkg/services"
)

func (a *API) GetServiceByName(w http.ResponseWriter, r *http.Request) {
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

	svc, err := a.svcService.GetServiceByName(r.Context(), user, name)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, svc)
}

func (a *API) GetServiceByID(w http.ResponseWriter, r *http.Request) {
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

	svc, err := a.svcService.GetServiceByID(r.Context(), user, id)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, svc)
}

func (a *API) CreateService(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var svc services.Service
	err = decoder.Decode(&svc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		a.sendErr(w, err)
		return
	}

	svc, err = a.svcService.CreateService(r.Context(), user, svc)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	a.sendJson(w, svc)
}

func (a *API) ListServices(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUser(r)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	svcirons, err := a.svcService.ListServices(r.Context(), user)
	if err != nil {
		a.sendErr(w, err)
		return
	}

	a.sendJson(w, svcirons)
}
