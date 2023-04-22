package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// home is a handler function
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// if we want to show 404 error for undefined paths, we can add this condition to the catch all case.

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	// used to read the template file into a template set
	// ParseFiles accepts variadic parametes, here files are unrolled in the function
	// execution
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "internal Server Error", 500)
		return
	}
	// execution of the template set allows us to write the template
	// into the response body
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internsal Server Error", 500)
	}

}

// snippetView is a handler function
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display specific snippet...."))
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.Header().Set("Allow", "Only POST")
		http.Error(w, "Method not Allowed", 405)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet id: %d", id)
}
