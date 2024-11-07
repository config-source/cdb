package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	if r.Status != 0 {
		return
	}

	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *StatusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func AccessLog(log zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			wr := &StatusRecorder{ResponseWriter: w, Status: 200}

			startTime := time.Now()
			next.ServeHTTP(wr, r)
			responseTime := time.Since(startTime)

			accessLog := log.Info()
			if wr.Status >= 400 {
				accessLog = log.Error()
			}

			accessLog.
				Str("url", r.URL.String()).
				Str("method", r.Method).
				Int("statusCode", wr.Status).
				Dur("responseTimeMilliseconds", responseTime).
				Msg("request served")
		},
	)

}
