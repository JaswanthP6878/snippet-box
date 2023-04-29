package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.jaswanthp.com/internal/models"
)

// home is a handler function
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// if we want to show 404 error for undefined paths, we can add this condition to the catch all case.

	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	// panic("oops fake panic heheh")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
	}

	app.render(w, http.StatusOK, "home.tmpl.html", &templateData{Snippets: snippets})

}

// snippetView is a handler function
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

	if r.Method != "POST" {
		w.Header().Set("Allow", "Only POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := "O snail"
	content := "this is a snail\nthis is what we do\ndo do be do\n"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
