package main

import (
	"log"
	"net/http"
)

func main() {
	// http.NewServeMux() inicializa um novo servemux
	mux := http.NewServeMux()
	// Registra home function como handler para "/" URL
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Cria um file server de "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// http.ListenAndServe() function inicia um new web server. Recebe 2 parâmetros:
	// TCP network address to listen on (no caso ":4000"),e o servemux
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
