package main

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func (app *Server) AddMidelware() {
	app.Router.Use(middleware.Heartbeat("/ping"))
	app.Router.Use(middleware.Recoverer)
	app.Router.Use(app.SessionLoad)
}

func (app *Server) SessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}
