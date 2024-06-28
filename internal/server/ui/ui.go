package ui

import (
	"html/template"
	"io"
	"net/http"

	"github.com/config-source/cdb/internal/configvalues"
	"github.com/config-source/cdb/internal/repository"
	"github.com/rs/zerolog"
)

type UI struct {
	repo               repository.ModelRepository
	configValueService *configvalues.Service
	log                zerolog.Logger
	templates          *template.Template
}

func New(
	repo repository.ModelRepository,
	configValueService *configvalues.Service,
	log zerolog.Logger,
	mux *http.ServeMux,
) *UI {
	ui := &UI{
		repo:               repo,
		configValueService: configValueService,
		log:                log,
		templates:          newTemplate(),
	}

	mux.HandleFunc("GET /", ui.Index)

	return ui
}

func (ui *UI) render(writer io.Writer, name string, data interface{}) {
	err := ui.templates.ExecuteTemplate(writer, name, nil)
	if err != nil {
		ui.log.Err(err).
			Str("templateName", name).
			Msg("failed to render template")
	}
}

// isHTMX checks if the Request was sent by HTMX.
func isHTMX(r *http.Request) bool {
	value := r.Header.Get("hx-request")
	return value != ""
}

func (ui *UI) Index(wr http.ResponseWriter, r *http.Request) {
	ui.render(wr, "pages/index", nil)
}
