package main

import (
	"errors"
	"fmt"
	"github.com/gitalek/gistpaste/pkg/helpers"
	"github.com/gitalek/gistpaste/pkg/models"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	gists, err := app.gists.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}


	data := &templateData{Gists: gists}
	app.render(w, r, "home.page.tmpl", data)
}

func (app *application) showGist(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
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

func (app *application) createGist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"

	id, err := app.gists.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect user to the relevant page for newly created gist.
	http.Redirect(w, r, fmt.Sprintf("/gist?id=%d", id), http.StatusSeeOther)
}
