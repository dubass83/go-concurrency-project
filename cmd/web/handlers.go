package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	data "github.com/dubass83/go-concurrency-project/data/sqlc"
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
		// send messages asynchronously with channels
		msg := Message{
			To:      []string{email},
			Subject: "Failed log in attempt",
			Data:    "invalid login attempt!",
		}
		app.Mail.MailerChan <- msg
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
	err := r.ParseForm()
	if err != nil {
		log.Error().Err(err).Msg("failed to parse the form from the request")
	}

	HashPass, err := utils.HashPassword(r.Form.Get("password"))
	if err != nil {
		log.Error().Err(err).Msg("failed to generate hash for a new user password from the request")
	}

	arg := data.InsertUserParams{
		Email: pgtype.Text{
			String: r.Form.Get("email"),
			Valid:  true,
		},
		FirstName: pgtype.Text{
			String: r.Form.Get("first-name"),
			Valid:  true,
		},
		LastName: pgtype.Text{
			String: r.Form.Get("last-name"),
			Valid:  true,
		},
		Password: pgtype.Text{
			String: HashPass,
			Valid:  true,
		},
		UserActive: pgtype.Int4{
			Int32: 0,
			Valid: true,
		},
	}
	_, err = app.Store.InsertUser(context.Background(), arg)
	if err != nil {
		log.Error().Err(err).Msg("failed insert a user to the database")
		app.Session.Put(r.Context(), "error", "Unable to create user.")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
	}

	// send activation email
	url := fmt.Sprintf("http://localhost:%s/activate?email=%s", app.Config.WebPort, r.Form.Get("email"))
	signedURL := app.GenerateTokenFromString(url)
	log.Info().Msg(signedURL)

	msg := Message{
		To:       []string{r.Form.Get("email")},
		Subject:  "Activate your account",
		Template: "confirmation-email",
		Data:     template.HTML(signedURL),
	}

	app.Mail.MailerChan <- msg

	// redirect user to the login page
	app.Session.Put(r.Context(), "flash", "Confirmation email sent. Check your email.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Server) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	// verify token from URL
	url := r.RequestURI
	testURL := fmt.Sprintf("http://localhost:%s%s", app.Config.WebPort, url)
	log.Info().Msg(testURL)
	ok := app.VerifyToken(testURL)

	if !ok {
		app.Session.Put(r.Context(), "error", "Invalid token.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Make user Active
	argEmaill := pgtype.Text{
		String: r.URL.Query().Get("email"),
		Valid:  true,
	}

	u, err := app.Store.GetUserByEmail(context.Background(), argEmaill)
	if err != nil {
		app.Session.Put(r.Context(), "error", "No such user in the Database.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	argUpdate := data.UpdateUserParams{
		ID: u.ID,
		UserActive: pgtype.Int4{
			Int32: 1,
			Valid: true,
		},
	}

	_, err = app.Store.UpdateUser(context.Background(), argUpdate)
	if err != nil {
		log.Error().Err(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "flash", "Acount is active. Now you can login to your account.")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
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

func (app *Server) ChooseSubscription(w http.ResponseWriter, r *http.Request) {
	if !app.Session.Exists(r.Context(), "userID") {
		app.Session.Put(r.Context(), "warning", "You must log in to see this page!")
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	arg := data.GetAllPlansParams{
		Limit:  10,
		Offset: 0,
	}
	plans, err := app.Store.GetAllPlans(context.Background(), arg)
	if err != nil {
		log.Error().Err(err)
		return
	}
	dataMap := make(map[string]any)
	dataMap["plans"] = plans

	app.render(w, r, "plans.page.gohtml", &TemplateData{
		DataMap: dataMap,
	})

}
