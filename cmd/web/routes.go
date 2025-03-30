package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *Server) MountHandlers() {
	app.Router.Get("/", app.HomePage)

	app.Router.Get("/login", app.LoginPage)
	app.Router.Post("/login", app.PostLoginPage)
	app.Router.Get("/logout", app.Logout)
	app.Router.Get("/register", app.RegisterPage)
	app.Router.Post("/register", app.PostRegisterPage)
	app.Router.Get("/activate", app.ActivateAccount)
	app.Router.Get("/test-email", app.SendTestEmail)

	app.Router.Mount("/members", app.AuthRouter())
}

func (app *Server) AuthRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Use(app.Auth)

	mux.Get("/plans", app.ChooseSubscription)
	mux.Get("/subscribe", app.SubscribeToPlan)

	return mux
}
