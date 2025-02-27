package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
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
	// User *data.User
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

	if td == nil {
		td = &TemplateData{}
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		app.ErrorLog.Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, app.AddDefaultData(td, r)); err != nil {
		app.ErrorLog.Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *Server) AddDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Authenticated = app.IsAuthenticated(r)
	td.Now = time.Now()

	return td
}

func (app *Server) IsAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "userID")
}
