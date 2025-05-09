package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	data "github.com/dubass83/go-concurrency-project/data/sqlc"
	"github.com/dubass83/go-concurrency-project/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
	"github.com/rs/zerolog/log"
)

type formatedPlans struct {
	ID                  int32     `json:"id"`
	PlanName            string    `json:"plan_name"`
	PlanAmount          int32     `json:"plan_amount"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	PlanAmountFormatted string    `json:"plan_amount_formatted"`
}

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
	log.Debug().Msgf("pass: %s  hash: %s", password, user.Password.String)
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

	//check if user has a plan add to the session
	userPlan, err := app.Store.GetOneUserPlan(context.Background(), pgtype.Int4{
		Int32: user.ID,
		Valid: true,
	})
	if err == nil {
		app.Session.Put(r.Context(), "user-plan", userPlan)
		// log.Debug().Any("user-plan", userPlan)
	}

	// login user
	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)

	app.Session.Put(r.Context(), "flash", "Successful login!")
	// redirect to the home page after successful login
	http.Redirect(w, r, "/", http.StatusSeeOther)
	log.Info().Msg("successfull log in")
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

	arg := data.GetAllPlansParams{
		Limit:  10,
		Offset: 0,
	}
	plans, err := app.Store.GetAllPlans(context.Background(), arg)
	if err != nil {
		log.Error().Err(err)
		return
	}

	fmtPlans := planAmountFormatted(plans)

	dataMap := make(map[string]any)
	dataMap["plans"] = fmtPlans

	app.render(w, r, "plans.page.gohtml", &TemplateData{
		DataMap: dataMap,
	})

}

func planAmountFormatted(plans []data.Plan) []formatedPlans {
	result := []formatedPlans{}
	for _, p := range plans {
		formatedPlan := formatedPlans{
			PlanAmountFormatted: fmt.Sprintf("$%.2f", float64(p.PlanAmount.Int32)/100),
			ID:                  p.ID,
			PlanName:            p.PlanName.String,
			PlanAmount:          p.PlanAmount.Int32,
			CreatedAt:           p.CreatedAt.Time,
			UpdatedAt:           p.UpdatedAt.Time,
		}
		result = append(result, formatedPlan)
	}
	return result
}

func (app *Server) SubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	planId, _ := strconv.Atoi(id)

	plan, err := app.Store.GetOnePlan(context.TODO(), int32(planId))
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to find plan!")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}

	user, ok := app.Session.Get(r.Context(), "user").(data.User)
	if !ok {
		app.Session.Put(r.Context(), "error", "Log in first!")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	app.Wait.Add(1)

	go func() {
		defer app.Wait.Done()

		invoice, err := app.getInvoice(&plan)
		if err != nil {
			app.ErrChan <- err
		}

		msg := Message{
			To:       []string{user.Email.String},
			Subject:  "Yuor invoice",
			Template: "invoice",
			Data:     invoice,
		}

		app.Mail.MailerChan <- msg
	}()

	app.Wait.Add(1)
	go func() {
		defer app.Wait.Done()

		pdf := app.generateManual(user, &plan)
		err := pdf.OutputFileAndClose(fmt.Sprintf("./tmp/%d_user_manual.pdf", user.ID))
		if err != nil {
			app.ErrChan <- err
			return
		}

		msg := Message{
			To:      []string{user.Email.String},
			Subject: "Yuor manual",
			Data:    "Your user manual is attached",
			AttachmentMap: map[string]string{
				"Manual.pdf": fmt.Sprintf("./tmp/%d_user_manual.pdf", user.ID),
			},
		}
		app.Mail.MailerChan <- msg

	}()

	// subscribe user to the choosen plan
	arg := data.SubscribeUserToPlanParams{
		UserID: user.ID,
		PlanID: plan.ID,
	}
	result, err := app.Store.SubscribeUserToPlan(context.TODO(), arg)
	if err != nil {
		log.Error().Err(err).Msg("failed to subscribe user to plan")
		app.Session.Put(r.Context(), "error", "Unable to subscribe to the plan!")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "user-plan", result.UserPlan)

	app.Session.Put(r.Context(), "flash", "Subscribed!")
	http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
}

func (app *Server) generateManual(u data.User, p *data.Plan) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(10, 13, 10)

	importer := gofpdi.NewImporter()

	time.Sleep(5 * time.Second)

	t := importer.ImportPage(pdf, "./pdf/manual.pdf", 1, "/MediaBox")
	pdf.AddPage()

	importer.UseImportedTemplate(pdf, t, 0, 0, 215.9, 0)

	pdf.SetX(75)
	pdf.SetY(150)

	pdf.SetFont("Arial", "", 12)

	pdf.MultiCell(0, 4, fmt.Sprintf("%s %s", u.FirstName.String, u.LastName.String), "", "C", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 4, fmt.Sprintf("%s User Guide", p.PlanName.String), "", "C", false)

	return pdf
}

func (app *Server) getInvoice(p *data.Plan) (string, error) {
	//dummy function
	return fmt.Sprintf("$%.2f", float64(p.PlanAmount.Int32)/100), nil
}
