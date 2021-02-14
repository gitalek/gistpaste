package main

import (
	"fmt"
	"github.com/gitalek/gistpaste/pkg/helpers"
	"html/template"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	filepaths := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(filepaths...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) showGist(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query().Get("id")
	id, err := helpers.ValidateParamId(p)
	if err != nil {
		w.WriteHeader(err.StatusCode)
		w.Write([]byte(err.Error()))
		return
	}
	fmt.Fprintf(w, "Display a specific gist with ID %d", id)
}

func (app *application) createGist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		statusCode := 405
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}
	w.Write([]byte("Create a new gist..."))
}
