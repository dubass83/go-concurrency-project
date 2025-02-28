package main

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func (app *Server) HomePage(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("start rendering the home page")
	app.render(w, r, "home.page.gohtml", nil)
}
