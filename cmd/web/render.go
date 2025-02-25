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

	if err := tmpl.Execute(w, nil); err != nil {
		app.ErrorLog.Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
