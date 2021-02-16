package main

import (
	"github.com/bmizerany/pat"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/gist/create", http.HandlerFunc(app.createGistForm))
	mux.Post("/gist/create", http.HandlerFunc(app.createGist))
	mux.Get("/gist/:id", http.HandlerFunc(app.showGist))

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
