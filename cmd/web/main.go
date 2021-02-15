package main

import (
	"database/sql"
	"flag"
	"github.com/gitalek/gistpaste/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

// application holds the application-wide deps
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	gists    *mysql.GistModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	// dsn is a data source name (or connection string)
	dsn := flag.String(
		"dsn",
		"web:password@/gistpaste?parseTime=true",
		"MySQL data source name",
	)
	flag.Parse()

	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infolog,
		gists:    &mysql.GistModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infolog.Printf("Starting server on %s\n", *addr)
	errorLog.Fatal(srv.ListenAndServe())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
