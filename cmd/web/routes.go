package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"snippetbox.jaswanthp.com/ui"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()
	// fileServer := http.FileServer(http.Dir("./ui/static"))
	fileServer := http.FileServer(http.FS(ui.Files))

	// we strip the /static before we reach the file handler becauase
	// if we keep it then it searches in ./ui/static/static which in not present

	// notfound errors are all served here

	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/Ping", Ping)

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))

	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))

	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignUp))

	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignUpPost))

	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))

	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	// protected routes for auth

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))

	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))

	// standard creates a chain of middlewares
	standard := alice.New(app.recoverPanic, app.logging, secureHeaders)
	return standard.Then(router)
}
