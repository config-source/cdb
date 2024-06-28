package ui

import (
	"embed"
	"html/template"
)

//go:embed templates
var templateFS embed.FS

func newTemplate() *template.Template {
	return template.Must(
		template.New("").ParseFS(
			templateFS,
			// this must be maintained such that each level of nesting used in
			// templates is represented here.
			"templates/*/*.html",
		),
	)
}
