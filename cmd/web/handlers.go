package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.jaswanthp.com/internal/models"
)

// home is a handler function
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// if we want to show 404 error for undefined paths, we can add this condition to the catch all case.
	// panic("oops fake panic heheh")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
	}

	app.render(w, http.StatusOK, "home.tmpl.html", &templateData{Snippets: snippets})

}

// snippetView is a handler function
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	// updated to getting params from the router context
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.render(w, http.StatusOK, "view.tmpl.html", &templateData{Snippet: snippet})

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "create.tmpl.html", data)

}
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
