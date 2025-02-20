package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/dubass83/go-concurrency-project/utils"
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

func initSessions(config utils.Config) *scs.SessionManager {
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
