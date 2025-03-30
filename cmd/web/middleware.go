package main

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func (app *Server) AddMidelware() {
	app.Router.Use(middleware.Logger)
	app.Router.Use(middleware.Heartbeat("/ping"))
	app.Router.Use(middleware.Recoverer)
	app.Router.Use(app.SessionLoad)
}

func (app *Server) SessionLoad(next http.Handler) http.Handler {
	log.Info().Msg("starting load and save session")
	return app.Session.LoadAndSave(next)
}

func (app *Server) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.Session.Exists(r.Context(), "userID") {
			app.Session.Put(r.Context(), "error", "Log in first!")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}
