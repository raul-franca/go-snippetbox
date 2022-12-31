package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {

	// Solução basica do go
	//mux := http.NewServeMux()
	//
	//mux.HandleFunc("/", app.home)
	//mux.HandleFunc("/snippet", app.showSnippet)
	//mux.HandleFunc("/snippet/create", app.createSnippet)
	//
	//// Cria um file server de "./ui/static" directory.
	//fileServer := http.FileServer(http.Dir("./ui/static/"))
	//mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	//return mux

	// Create a new middleware chain containing the middleware specific to
	//our dynamic application routes. For now, this chain will only contain
	//the session middleware but we'll add more to it later.
	dynamicMiddleware := alice.New(app.session.Enable)

	// com alice -> http dMiddleware e bmizerany/pat
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))
	// User routes.
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	// Leave the static files route unchanged.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
