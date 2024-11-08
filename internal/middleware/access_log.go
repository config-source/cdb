package middleware

import (
	"net/http"
	"time"

	"github.com/config-source/cdb/internal/apiutils"
	"github.com/rs/zerolog"
)

func AccessLog(log zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			wr := apiutils.NewStatusRecorder(w)

			startTime := time.Now()
			next.ServeHTTP(wr, r)
			responseTime := time.Since(startTime)

			accessLog := log.Info()
			if wr.Status() >= 400 {
				accessLog = log.Error()
			}

			accessLog.
				Str("url", r.URL.String()).
				Str("method", r.Method).
				Int("statusCode", wr.Status()).
				Dur("responseTimeMilliseconds", responseTime).
				Msg("request served")
		},
	)

}
