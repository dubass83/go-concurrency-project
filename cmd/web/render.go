package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	data "github.com/dubass83/go-concurrency-project/data/sqlc"
	"github.com/rs/zerolog/log"
)

type TemplateData struct {
	StringMap     map[string]string
	IntMap        map[string]int
	FloatMap      map[string]float64
	DataMap       map[string]any
	Flash         string
	Warning       string
	Error         string
	Authenticated bool
	Now           time.Time
	User          *data.User
	UPlan         *data.UserPlan
}

func (app *Server) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) {
	partials := []string{
		fmt.Sprintf("%s/base.layout.gohtml", app.Config.PathToTemplate),
		fmt.Sprintf("%s/header.partial.gohtml", app.Config.PathToTemplate),
		fmt.Sprintf("%s/navbar.partial.gohtml", app.Config.PathToTemplate),
		fmt.Sprintf("%s/footer.partial.gohtml", app.Config.PathToTemplate),
		fmt.Sprintf("%s/alerts.partial.gohtml", app.Config.PathToTemplate),
	}

	var templateSlice []string
	templateSlice = append(templateSlice, partials...)
	templateSlice = append(templateSlice, fmt.Sprintf("%s/%s", app.Config.PathToTemplate, t))
	// log.Debug().Msgf("templateSlice: %v", templateSlice)
	// log.Debug().Msgf("Template path: %s", app.Config.PathToTemplate)

	// Check if files exist
	for _, path := range templateSlice {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Error().Msgf("Template file does not exist: %s", path)
			http.Error(w, "Template file not found", http.StatusInternalServerError)
			return
		}
	}
	if td == nil {
		td = &TemplateData{}
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse files")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", app.AddDefaultData(td, r)); err != nil {
		log.Error().Err(err).Msg("failed to execute template with data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *Server) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	if app.IsAuthenticated(r) {
		td.Authenticated = true
		if u, ok := app.Session.Get(r.Context(), "user").(data.User); ok {
			td.User = &u
		}
		if up, ok := app.Session.Get(r.Context(), "user-plan").(data.UserPlan); ok {
			td.UPlan = &up
		}
	}
	td.Now = time.Now()
	// log.Debug().Msgf("td: %v", td)
	return td
}

func (app *Server) IsAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "userID")
}
