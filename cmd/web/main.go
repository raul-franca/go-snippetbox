package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// http.NewServeMux() inicializa um novo servemux
	mux := http.NewServeMux()

	// Define uma flag de linha de comando com o nome 'addr', um valor padrão de ":4000"
	//e um texto ajuda explicando o que o flag controla
	addr := flag.String("addr", ":4000", "HTTP network address")
	//A função flag.Parse() le se a flag foi usada no command-line e altera seu valor se usada
	//no path root $ go run . -addr=":8000"
	flag.Parse()

	// Registra home function como handler para "/" URL
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Cria um file server de "./ui/static" directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// http.ListenAndServe() function inicia um new web server. Recebe 2 parâmetros:
	// TCP network address to listen on (no caso ":4000"),e o servemux
	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
