package main

import (
	"errors"
	"fmt"
	"github.com/gitalek/gistpaste/pkg/forms"
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

func (app *application) createGistForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createGist(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
	}

	id, err := app.gists.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Gist successfully created!")

	// Redirect user to the relevant page for newly created gist.
	http.Redirect(w, r, fmt.Sprintf("/gist/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm (w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Display the user signup form...")
}

func (app *application) signupUser (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Create a new user...")
}

func (app *application) loginUserForm (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Display the user login form...")
}

func (app *application) loginUser (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Authenticate and login the user...")
}

func (app *application) logoutUser (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Logout the user...")
}
