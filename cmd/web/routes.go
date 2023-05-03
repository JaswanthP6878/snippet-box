package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// we strip the /static before we reach the file handler becauase
	// if we keep it then it searches in ./ui/static/static which in not present

	// notfound errors are all served here

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))

	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))

	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))

	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignUp))

	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignUpPost))

	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))

	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))

	// standard creates a chain of middlewares
	standard := alice.New(app.recoverPanic, app.logging, secureHeaders)
	return standard.Then(router)
}
