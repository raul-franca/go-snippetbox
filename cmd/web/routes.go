package main

import "net/http"

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Cria um file server de "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	//return app.logRequest(secureHeaders(mux))
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
