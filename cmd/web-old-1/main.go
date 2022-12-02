package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/raul-franca/go-snippetbox/pkg/models/mysql"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
}

func main() {

	// Define uma flag de linha de comando com o nome 'addr', um valor padrão de ":4000"
	//e um texto ajuda explicando o que o flag controla
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Define uma nova command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "root:root@/snippetbox?parseTime=true", "MySQL data source name")

	//A função flag.Parse() le se a flag foi usada no command-line e altera seu valor se usada
	//no path root $ go run . -addr=":8000"
	flag.Parse()

	//log.New() criar um logger para escrever mensagens de informação.
	//Parâmetros: o destino para gravar os logs (os.Stdout), uma string
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{DB: db},
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
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}

func openDB(dsn string) (*sql.DB, error) {
	// sql.Open() function initializes a new sql.DB object, which is essentially a
	//pool of database connections.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
