package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	data "github.com/dubass83/go-concurrency-project/data/sqlc"
	"github.com/dubass83/go-concurrency-project/utils"
	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// read config
	conf, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot load configuration")
	}
	if conf.Enviroment == "devel" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		// log.Debug().Msgf("config values: %+v", conf)
	}

	// ctx, stop := signal.NotifyContext(context.Background(), interaptSignals...)
	// defer stop()
	// connect to the database
	connPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig(conf))
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot validate the db connection string")
	}
	err = connPool.Ping(context.TODO())
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot ping the database from the connection pool")
	}

	// create sessions
	session := initSessions(conf)

	// create channels

	// create waitgroup
	wg := sync.WaitGroup{}

	// set up the application config
	app := Server{
		Config:  conf,
		Router:  chi.NewRouter(),
		Store:   data.NewStore(connPool),
		Session: session,
		Wait:    &wg,
	}
	// run db migration
	app.runDbMigration()

	// set up mail

	// listen for the signals
	go app.ListenForShutdown()

	// listen for web connections
	app.serve()
}

func (app *Server) serve() {
	app.AddMidelware()
	app.MountHandlers()
	// start http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.Config.WebPort),
		Handler: app.Router,
	}

	log.Info().Msgf("starting http web server on the port: %s", app.Config.WebPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to start http web server on the port: %s", app.Config.WebPort)
	}
}

// PoolConfig create config for db connection pool
func poolConfig(conf utils.Config) *pgxpool.Config {

	dbConfig, err := pgxpool.ParseConfig(conf.DBSource)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to create a config")
	}
	if conf.DBPoolMaxConns != 0 {
		dbConfig.MaxConns = conf.DBPoolMaxConns
	}
	if conf.DBPoolMinConns != 0 {
		dbConfig.MinConns = conf.DBPoolMinConns
	}
	if conf.DBPoolMaxConnLifetime != time.Second*0 {
		dbConfig.MaxConnLifetime = conf.DBPoolMaxConnLifetime
	}
	if conf.DBPoolMaxConnIdleTime != time.Second*0 {
		dbConfig.MaxConnIdleTime = conf.DBPoolMaxConnIdleTime
	}
	if conf.DBPoolHealthCheckPeriod != time.Second*0 {
		dbConfig.HealthCheckPeriod = conf.DBPoolHealthCheckPeriod
	}
	if conf.DBPoolConnectTimeout != time.Second*0 {
		dbConfig.ConnConfig.ConnectTimeout = conf.DBPoolConnectTimeout
	}

	return dbConfig
}

func initSessions(config utils.Config) *scs.SessionManager {
	gob.Register(data.User{})
	// setup session
	session := scs.New()
	session.Store = redisstore.New(initRedis(config))
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session
}

func initRedis(config utils.Config) *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", config.RedisURL)
		},
	}
	return redisPool
}

func (app *Server) ListenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutdown()
	os.Exit(0)
}

func (app *Server) shutdown() {
	log.Info().Msg("starting shutdown process for the app...")

	app.Wait.Wait()

	log.Info().Msg("all chanels will be stoped and app will be prepared for gracefully shutdown")
	// TODO close all chanels
}

// runDbMigration run db migration from file to db
func (app *Server) runDbMigration() {
	m, err := migrate.New(app.Config.MigrationURL, app.Config.DBSource)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("can not create migration instance")
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().
			Err(err).
			Msg("can not run migration up")
	}
	log.Info().Msg("successfully run db migration")
}
