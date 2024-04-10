package middleware

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"time"

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

func (r *StatusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("http: proxy error: can't switch protocols using non-Hijacker ResponseWriter type")
	}

	return h.Hijack()
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
