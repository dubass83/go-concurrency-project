package main

import (
	"context"
	"os"
	"time"

	"github.com/dubass83/go-concurrency-project/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// connect to the database
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

	connPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig(conf))
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot validate db connection string")
	}
	err = connPool.Ping(context.TODO())
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot ping database")
	}

	// create sessions

	// create channels

	// create waitgroup

	// set up the application config

	// set up mail

	// listen for web connections
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
