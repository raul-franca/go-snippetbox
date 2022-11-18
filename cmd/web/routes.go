package main

import "net/http"

func (app application) routes() *http.ServeMux {

	// http.NewServeMux() inicializa um novo servemux
	mux := http.NewServeMux()
	// Registra home function como handler para "/" URL
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// Cria um file server de "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
