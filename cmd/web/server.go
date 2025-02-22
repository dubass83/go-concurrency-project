package main

import (
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/dubass83/go-concurrency-project/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Server struct {
	Config   utils.Config
	Router   *chi.Mux
	Session  *scs.SessionManager
	Db       *pgxpool.Pool
	InfoLog  *zerolog.Event
	ErrorLog *zerolog.Event
	Wait     *sync.WaitGroup
}

func (app *Server) AddMidelware() {
	app.Router.Use(middleware.Heartbeat("/ping"))
	app.Router.Use(middleware.Recoverer)
}
