package main

import (
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/dubass83/go-concurrency-project/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Config  utils.Config
	Router  *chi.Mux
	Session *scs.SessionManager
	Db      *pgxpool.Pool
	Wait    *sync.WaitGroup
}
