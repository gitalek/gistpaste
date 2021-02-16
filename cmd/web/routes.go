package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/gist", app.showGist)
	mux.HandleFunc("/gist/create", app.createGist)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mws := []middleware{
		secureHeaders,
		app.logRequest,
		app.recoverPanic,
	}
	return wrapMux(mux, mws...)
}

func wrapMux(mux *http.ServeMux, mws ...middleware) http.Handler {
	var h http.Handler
	h = mux
	for _, mw := range mws {
		h = mw(h)
	}
	return h
}
