package main

import (
	"sync"

	"github.com/alexedwards/scs/v2"
	data "github.com/dubass83/go-concurrency-project/data/sqlc"
	"github.com/dubass83/go-concurrency-project/utils"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Config  utils.Config
	Router  *chi.Mux
	Session *scs.SessionManager
	Store   data.Store
	Wait    *sync.WaitGroup
	Mail    Mail
}
