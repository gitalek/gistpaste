package main

import (
	"github.com/bmizerany/pat"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := pat.New()
	mux.Get("/", app.session.Enable(http.HandlerFunc(app.home)))
	mux.Get("/gist/create", app.session.Enable(http.HandlerFunc(app.createGistForm)))
	mux.Post("/gist/create", app.session.Enable(http.HandlerFunc(app.createGist)))
	mux.Get("/gist/:id", app.session.Enable(http.HandlerFunc(app.showGist)))

	mux.Get("/user/signup", app.session.Enable(http.HandlerFunc(app.signupUserForm)))
	mux.Post("/user/signup", app.session.Enable(http.HandlerFunc(app.signupUser)))
	mux.Get("/user/login", app.session.Enable(http.HandlerFunc(app.loginUserForm)))
	mux.Post("/user/login", app.session.Enable(http.HandlerFunc(app.loginUser)))
	mux.Post("/user/logout", app.session.Enable(http.HandlerFunc(app.logoutUser)))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	mws := []middleware{
		secureHeaders,
		app.logRequest,
		app.recoverPanic,
	}
	return wrapMux(mux, mws...)
}

func wrapMux(mux *pat.PatternServeMux, mws ...middleware) http.Handler {
	var h http.Handler
	h = mux
	for _, mw := range mws {
		h = mw(h)
	}
	return h
}
