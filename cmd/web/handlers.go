package main

import (
	"errors"
	"fmt"
	"github.com/gitalek/gistpaste/pkg/helpers"
	"github.com/gitalek/gistpaste/pkg/models"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	gists, err := app.gists.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}


	data := &templateData{Gists: gists}
	app.render(w, r, "home.page.tmpl", data)
}

func (app *application) showGist(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(":id")
	// todo: app method for ValidateParamIde
	errValid := helpers.ValidateParamId(id)
	if errValid != nil {
		w.WriteHeader(errValid.StatusCode)
		w.Write([]byte(errValid.Error()))
		return
	}

	g, err := app.gists.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := &templateData{Gist: g}
	app.render(w, r, "show.page.tmpl", data)
}

func (app *application) createGistForm(w http.ResponseWriter, r *http.Request)  {
	app.render(w, r, "create.page.tmpl", nil)
}

func (app *application) createGist(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	id, err := app.gists.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect user to the relevant page for newly created gist.
	http.Redirect(w, r, fmt.Sprintf("/gist/%d", id), http.StatusSeeOther)
}
