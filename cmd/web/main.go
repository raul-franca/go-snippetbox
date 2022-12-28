package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/raul-franca/go-snippetbox/pkg/models/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	session       *sessions.Session
	templateCache map[string]*template.Template
}

func main() {

	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	// Define uma flag de linha de comando com o nome 'addr', um valor padrão de ":4000"
	//e um texto ajuda explicando o que o flag controla
	//no path root $ go run . -addr=":8000"
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Define uma nova command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "root:root@/snippetbox?parseTime=true", "MySQL data source name")

	//A função flag.Parse() le se a flag foi usada no command-line e altera seu valor se usada

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

	// Initialize a new template cache...
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	// Inicializar uma nova instance a application contendo as dependências.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// cria um struct http.Server com Addr, Handler, ErrorLog personalizado
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

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
