package main

func (app *Server) MountHandlers() {
	app.Router.Get("/", app.HomePage)
}
