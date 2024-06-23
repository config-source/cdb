package ui

import "net/http"

type UI struct {
}

func New(mux *http.ServeMux) *UI {
	return &UI{}
}

// isHTMX returns a boolean indicating that the given request was sent by HTMX.
func isHTMX(r *http.Request) bool {
	value := r.Header.Get("hx-request")
	return value != ""
}
