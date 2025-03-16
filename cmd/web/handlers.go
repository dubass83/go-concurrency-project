package main

import (
	"context"
	"net/http"

	"github.com/dubass83/go-concurrency-project/utils"
	"github.com/jackc/pgx/v5/pgtype"
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
	// renew session token every time when user post the login form
	_ = app.Session.RenewToken(r.Context())

	// parse the form
	err := r.ParseForm()
	if err != nil {
		log.Error().Err(err).Msg("failed to parse login form")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	getUserEmail := pgtype.Text{
		String: email,
		Valid:  true,
	}

	// get user from the database
	user, err := app.Store.GetUserByEmail(context.TODO(), getUserEmail)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user by email from database")
		app.Session.Put(r.Context(), "error", "invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// check user password
	err = utils.CheckPassword(password, user.Password.String)
	if err != nil {
		log.Error().Err(err).Msg("invalid credentials")
		app.Session.Put(r.Context(), "error", "invalid credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// login user
	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)

	app.Session.Put(r.Context(), "flash", "Successful login!")
	// redirect to the home page after successful login
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Server) Logout(w http.ResponseWriter, r *http.Request) {
	// cleanup the Session
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Server) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}

func (app *Server) PostRegisterPage(w http.ResponseWriter, r *http.Request) {

}

func (app *Server) ActivateAccount(w http.ResponseWriter, r *http.Request) {

}

func (app *Server) SendTestEmail(w http.ResponseWriter, r *http.Request) {

	email := Message{
		From:      "Maks",
		FromEmail: "dubass@test.work",
		Subject:   "test message",
		To:        []string{"boloto@test.com", "sasas@test.com"},
		Data:      "Hello world",
	}

	app.Mail.Sender.SendEmail(email, app.Mail.ErrChan)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
