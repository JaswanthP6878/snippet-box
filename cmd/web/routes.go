package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static"))
	// we strip the /static before we reach the file handler becauase
	// if we keep it then it searches in ./ui/static/static which in not present
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)
	// standard creates a chain of middlewares
	standard := alice.New(app.recoverPanic, app.logging, secureHeaders)
	return standard.Then(mux)
}
