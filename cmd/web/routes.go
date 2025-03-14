package main

func (app *Server) MountHandlers() {
	app.Router.Get("/", app.HomePage)

	app.Router.Get("/login", app.LoginPage)
	app.Router.Post("/login", app.PostLoginPage)
	app.Router.Get("/logout", app.Logout)
	app.Router.Get("/register", app.RegisterPage)
	app.Router.Post("/register", app.PostRegisterPage)
	app.Router.Get("/activate-account", app.ActivateAccount)

	app.Router.Get("/test-email", app.SendTestEmail)
}
