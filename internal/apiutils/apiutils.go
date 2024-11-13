package apiutils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/config-source/cdb/pkg/auth"
	"github.com/config-source/cdb/pkg/configkeys"
	"github.com/config-source/cdb/pkg/configvalues"
	"github.com/config-source/cdb/pkg/environments"
	"github.com/config-source/cdb/pkg/services"
	"github.com/rs/zerolog"
)

type StatusRecorder struct {
	http.ResponseWriter
	status int
}

func NewStatusRecorder(w http.ResponseWriter) *StatusRecorder {
	return &StatusRecorder{
		ResponseWriter: w,
	}
}

func (r *StatusRecorder) WriteHeader(status int) {
	if r.status != 0 {
		return
	}

	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *StatusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func (r *StatusRecorder) Status() int {
	if r.status == 0 {
		return 200
	}

	return r.status
}

type ErrorResponse struct {
	Message string
}

func (er ErrorResponse) Error() string {
	return er.Message
}

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{
		Message: msg,
	}
}

func SendJSON(log zerolog.Logger, w http.ResponseWriter, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Err(err).Msg("failed to encode a payload")
	}
}

func SendErr(log zerolog.Logger, w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case
		errors.Is(err, auth.ErrUserNotFound),
		errors.Is(err, environments.ErrNotFound),
		errors.Is(err, configkeys.ErrNotFound),
		errors.Is(err, services.ErrNotFound),
		errors.Is(err, configvalues.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case
		errors.Is(err, configvalues.ErrNotValid),
		errors.Is(err, configvalues.ErrAlreadySet),
		errors.Is(err, auth.ErrPublicRegisterDisabled),
		errors.Is(err, auth.ErrEmailInUse):
		w.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, auth.ErrUnauthorized),
		errors.Is(err, auth.ErrInvalidPassword):
		w.WriteHeader(http.StatusForbidden)
	case errors.Is(err, auth.ErrUnauthenticated):
		w.WriteHeader(http.StatusUnauthorized)
	// This is safe because subsequent calls to WriteHeader are ignored so
	// callers can set the status code before calling errorResponse but if they
	// haven't we want to send a 500.
	default:
		switch v := w.(type) {
		case *StatusRecorder:
			if v.Status() == 0 {
				log.Error().
					Err(err).
					Str("method", r.Method).
					Str("url", r.URL.Path).
					Msg("unhandled error")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}

	SendJSON(log, w, NewErrorResponse(err.Error()))
}
