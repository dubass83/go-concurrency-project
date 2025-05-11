package main

import (
	"encoding/gob"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	data "github.com/dubass83/go-concurrency-project/data/sqlc"
	"github.com/dubass83/go-concurrency-project/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var testApp Server

func TestMain(m *testing.M) {

	// read config
	config, err := utils.LoadConfig("../../test_conf")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot load test configuration")
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// session setup
	gob.Register(data.User{})
	gob.Register(data.UserPlan{})
	gob.Register(pgtype.Int4{})
	gob.Register(pgtype.Text{})
	gob.Register(pgtype.Timestamp{})

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	// create channels
	mailChan := make(chan Message, 100)
	errChan := make(chan error)
	doneChan := make(chan bool)

	// set up mail
	sender, err := NewMailSender(config)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("failed to create new mail sender")
	}
	mail := Mail{
		MailerChan: mailChan,
		ErrChan:    errChan,
		DoneChan:   doneChan,
		Sender:     sender,
	}

	testApp = Server{
		Config:      config,
		Router:      chi.NewRouter(),
		Store:       nil,
		Session:     session,
		Wait:        &sync.WaitGroup{},
		Mail:        mail,
		ErrChan:     make(chan error),
		ErrChanDone: make(chan bool),
	}

	go func() {
		for {
			select {
			case <-testApp.Mail.MailerChan:
				testApp.Wait.Done()
			case <-testApp.Mail.ErrChan:
			case <-testApp.Mail.DoneChan:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case err := <-testApp.ErrChan:
				log.Error().Err(err)
			case <-testApp.ErrChanDone:
				return
			}
		}
	}()

	testApp.AddMidelware()
	testApp.MountHandlers()

	os.Exit(m.Run())
}
