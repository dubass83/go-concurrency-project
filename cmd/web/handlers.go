package main

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func (app *Server) HomePage(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("start rendering the home page")
	app.render(w, r, "home.page.gohtml", nil)
}

func (app *Server) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}

func (app *Server) PostLoginPage(w http.ResponseWriter, r *http.Request) {

}

func (app *Server) Logout(w http.ResponseWriter, r *http.Request) {

}

func (app *Server) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}

func (app *Server) PostRegisterPage(w http.ResponseWriter, r *http.Request) {

}

func (app *Server) ActivateAccount(w http.ResponseWriter, r *http.Request) {

}
