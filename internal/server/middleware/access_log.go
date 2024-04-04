package middleware

import (
	"net/http"

	"github.com/rs/zerolog"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func AccessLog(log zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			wr := &StatusRecorder{ResponseWriter: w, Status: 200}
			wr.Header().Set("Content-Type", "application/json")

			next.ServeHTTP(wr, r)

			accessLog := log.Info()
			if wr.Status >= 400 {
				accessLog = log.Error()
			}
			accessLog.
				Str("url", r.URL.String()).
				Str("method", r.Method).
				Int("statusCode", wr.Status).
				Msg("request served")
		},
	)

}
