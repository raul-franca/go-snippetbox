package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	// Define uma flag de linha de comando com o nome 'addr', um valor padrão de ":4000"
	//e um texto ajuda explicando o que o flag controla
	addr := flag.String("addr", ":4000", "HTTP network address")
	//A função flag.Parse() le se a flag foi usada no command-line e altera seu valor se usada
	//no path root $ go run . -addr=":8000"
	flag.Parse()

	//log.New() criar um logger para escrever mensagens de informação.
	//Parâmetros: o destino para gravar os logs (os.Stdout), uma string
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// cria um struct http.Server com Addr, Handler, ErrorLog personalizado
	// aumenta o codigo mas facilita o entendimento
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	//log.Printf("Starting server on %s", *addr)
	infoLog.Printf("Starting server on %s", *addr)

	//err := http.ListenAndServe(*addr, mux)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)

}
