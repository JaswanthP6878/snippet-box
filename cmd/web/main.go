package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"snippetbox.jaswanthp.com/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
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

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	// closes the connection pool when graceful termination happens
	defer db.Close()
	mux := app.routes()

	sv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = sv.ListenAndServe()
	errorLog.Fatal(err)
}

// function to open connetion pool for Mysql database for a given dsn
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
