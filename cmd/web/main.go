package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"snippetbox.jaswanthp.com/internal/models"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// declaring a flag variable, with default value ":4000" and
	// flag name as "addr" and a message
	addr := flag.String("addr", ":4000", "Defining the port address")
	// dsn is a string to identify the connection to mysql. it defines the data source.
	dsn := flag.String("dsn", "bill:passpass@/snippetbox?parseTime=true", "MYSQL data source name")

	// used to pasrse the command line value and get flags and assign
	// them to declared flags.
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error:\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := OpenDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	// closes the connection pool when graceful termination happens
	defer db.Close()
	mux := app.routes()

	sv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   mux,
		TLSConfig: tlsConfig,

		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	// HTTPS server
	err = sv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

// function to open connection pool for Mysql database for a given dsn
func OpenDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}
